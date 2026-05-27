#include <windows.h>
#include <stdio.h>
#include <winsock.h>

#define DEFAULT_PORT 8888
#define BUFFER_SIZE 256

int main(int argc, char** argv) {
    WSADATA WSAData;
    WORD wWSAVer = MAKEWORD(2, 2);
    SOCKET sock;
    struct sockaddr_in server_addr;
    char buffer[BUFFER_SIZE];
    int server_addr_size = sizeof(server_addr);
    char server_ip[50];

    printf("=== UDP Time Client ===\n");
    printf("Enter server IP address: ");
    fgets(server_ip, sizeof(server_ip), stdin);
    server_ip[strcspn(server_ip, "\n")] = 0;

    // Инициализация Winsock
    if (WSAStartup(wWSAVer, &WSAData) != 0) {
        printf("WSAStartup failed. Error Code: %d\n", WSAGetLastError());
        return -1;
    }

    // Создание сокета
    sock = socket(AF_INET, SOCK_DGRAM, 0);
    if (sock == INVALID_SOCKET) {
        printf("Socket creation failed. Error Code: %d\n", WSAGetLastError());
        WSACleanup();
        return -1;
    }

    // Настройка адреса сервера
    memset(&server_addr, 0, sizeof(server_addr));
    server_addr.sin_family = AF_INET;
    server_addr.sin_addr.s_addr = inet_addr(server_ip);
    server_addr.sin_port = htons(DEFAULT_PORT);

    if (server_addr.sin_addr.s_addr == INADDR_NONE) {
        printf("Invalid IP address!\n");
        closesocket(sock);
        WSACleanup();
        return -1;
    }

    printf("Connected to time server at %s:%d\n", server_ip, DEFAULT_PORT);
    printf("Type 'exit' to quit.\n");
    printf("Enter time in format HH.MM.SS (e.g., 12.01.20):\n\n");

    while (1) {
        printf("> ");
        fgets(buffer, BUFFER_SIZE, stdin);
        buffer[strcspn(buffer, "\n")] = 0;
        
        if (strcmp(buffer, "exit") == 0) break;
        
        // Отправка строки времени
        if (sendto(sock, buffer, strlen(buffer), 0,
                   (struct sockaddr*)&server_addr, sizeof(server_addr)) == SOCKET_ERROR) {
            printf("Send failed. Error Code: %d\n", WSAGetLastError());
            continue;
        }
        
        // Получение ответа
        memset(buffer, 0, BUFFER_SIZE);
        int bytes_recv = recvfrom(sock, buffer, BUFFER_SIZE - 1, 0,
                                  (struct sockaddr*)&server_addr, &server_addr_size);
        
        if (bytes_recv > 0) {
            buffer[bytes_recv] = '\0';
            printf("Result: %s\n\n", buffer);
        } else {
            printf("Receive failed. Error Code: %d\n", WSAGetLastError());
        }
    }

    closesocket(sock);
    WSACleanup();
    return 0;
}