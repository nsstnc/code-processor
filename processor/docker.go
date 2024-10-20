package processor

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func RunDockerContainer(code string, lang string) (string, error) {
	var fileExtension string
	switch lang {
	case "python":
		fileExtension = ".py"
	case "cpp":
		fileExtension = ".cpp"
	case "c":
		fileExtension = ".c"
	default:
		return "", fmt.Errorf("unsupported language: %s", lang)
	}

	cli, err := client.NewClientWithOpts(client.WithVersion("1.45"))
	if err != nil {
		return "", err
	}

	// Создание временного файла с кодом
	tempFile, err := os.CreateTemp(os.TempDir(), "code-*"+fileExtension)
	if err != nil {
		return "", err
	}
	defer os.Remove(tempFile.Name()) // Удалить файл после использования

	// Запись кода в файл
	if _, err := tempFile.Write([]byte(code)); err != nil {
		return "", err
	}
	tempFile.Close()

	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}
	processorDir := filepath.Join(currentDir, "processor")

	// Сборка Docker-образа
	cmd := exec.Command("docker", "build", "-t", "temp_code_image", ".")
	cmd.Dir = processorDir
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("error building Docker image: %w", err)
	}

	_, fileName := filepath.Split(tempFile.Name())

	// Создание контейнера
	resp, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image: "temp_code_image",
		Cmd:   []string{"./run_code.sh", lang, "/usr/src/app/" + fileName},
	}, nil, nil, nil, "")
	if err != nil {
		return "", err
	}

	// Копируем временный файл в контейнер
	copyCmd := exec.Command("docker", "cp", tempFile.Name(), fmt.Sprintf("%s:/usr/src/app/%s", resp.ID, fileName))
	if err := copyCmd.Run(); err != nil {
		return "", fmt.Errorf("failed to copy file to container: %w", err)
	}

	// Запуск контейнера
	if err := cli.ContainerStart(context.Background(), resp.ID, container.StartOptions{}); err != nil {
		return "", err
	}

	// Ожидание завершения контейнера
	statusCh, errCh := cli.ContainerWait(context.Background(), resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return "", fmt.Errorf("error waiting for container: %w", err)
		}
	case <-statusCh:
		// Контейнер завершил работу
	}

	// Получение логов
	out, err := cli.ContainerLogs(context.Background(), resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return "", fmt.Errorf("failed to get container logs: %w", err)
	}
	defer out.Close()

	logs, err := io.ReadAll(out)
	if err != nil {
		return "", fmt.Errorf("failed to read logs: %w", err)
	}

	// Остановка и удаление контейнера
	if err := cli.ContainerRemove(context.Background(), resp.ID, container.RemoveOptions{Force: true}); err != nil {
		return "", fmt.Errorf("error removing Docker container: %w", err)
	}

	// Удаление Docker-образа
	cmd = exec.Command("docker", "rmi", "temp_code_image")
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("error removing Docker image: %w", err)
	}

	return string(logs), nil
}
