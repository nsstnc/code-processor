package processor

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func RunDockerContainer(code string, lang string) (string, error) {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return "", err
	}

	// Создание временного файла с кодом
	tempFile, err := os.CreateTemp(os.TempDir(), "code-*.txt")
	if err != nil {
		return "", err
	}
	defer os.Remove(tempFile.Name()) // Удалить файл после использования

	// Запись кода в файл
	if _, err := tempFile.Write([]byte(code)); err != nil {
		return "", err
	}
	tempFile.Close()

	// Создание контейнера
	resp, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image: "your_docker_image", // Имя вашего Docker образа
		Cmd:   []string{lang, tempFile.Name()},
	}, nil, nil, nil, "")
	if err != nil {
		return "", err
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

	return string(logs), nil
}
