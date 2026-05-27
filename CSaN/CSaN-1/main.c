#include <winsock2.h>
#include <stdio.h>
#include <string.h>
#include <ws2tcpip.h> 
#include "date.h"

int main(int argc, char** argv) {
    int use_udp = 0;
    if (argc > 1 && strcmp(argv[1], "-u") == 0) {
        use_udp = 1;
        printf("Using UDP protocol\n");
    } else {
        printf("Using TCP protocol (default)\n");
    }

    WSADATA wsa;
    if (WSAStartup(MAKEWORD(2, 2), &wsa) != 0) {
        printf("error: WSA failed\n");
        return 1;
    }

    int fd;
    if (use_udp) {
        fd = socket(AF_INET, SOCK_DGRAM, 0);
    } else {
        fd = socket(AF_INET, SOCK_STREAM, 0);
    }
    
    if (fd == INVALID_SOCKET) {
        printf("error: Socket creation failed: %d\n", WSAGetLastError());
        WSACleanup();
        return 1;
    }
    
    struct sockaddr_in addr;
    addr.sin_family = AF_INET;
    addr.sin_addr.s_addr = INADDR_ANY;
    addr.sin_port = htons(8080);

    if (bind(fd, (struct sockaddr*)&addr, sizeof(addr)) == SOCKET_ERROR) {
        printf("error: Bind failed: %d\n", WSAGetLastError());
        closesocket(fd);
        WSACleanup();
        return 1;
    }

    if (!use_udp) {
        if (listen(fd, 10) == SOCKET_ERROR) {
            printf("error: Listen failed: %d\n", WSAGetLastError());
            closesocket(fd);
            WSACleanup();
            return 1;
        }

        printf("TCP Server listening on port 8080...\n");
        
        struct sockaddr_in client_addr;
        socklen_t client_len = sizeof(client_addr);

        int client_fd = accept(fd, (struct sockaddr*)&client_addr, &client_len);
        
        if (client_fd == INVALID_SOCKET) {
            printf("error: Accept failed: %d\n", WSAGetLastError());
            closesocket(fd);
            WSACleanup();
            return 1;
        }
        
        printf("Client connected!\n");
        
        char first_date[1024];
        char second_date[1024];
        
        while (1) {
            memset(first_date, 0, sizeof(first_date));
            memset(second_date, 0, sizeof(second_date));
            
            int bytes_received = recv(client_fd, first_date, sizeof(first_date) - 1, 0);
            if (bytes_received <= 0) {
                printf("Client disconnected or error occurred\n");
                break;
            }
            
            bytes_received = recv(client_fd, second_date, sizeof(second_date) - 1, 0);
            if (bytes_received <= 0) {
                printf("Client disconnected or error occurred\n");
                break;
            }

            char res[1024];
            sprintf(res, "Count of days between two dates is %d", daysBetweenDates(first_date, second_date));
            printf("%s\n", res);
            send(client_fd, res, strlen(res) + 1, 0);
        }

        closesocket(client_fd);
    } else {
        printf("UDP Server listening on port 8080...\n");
        
        struct sockaddr_in client_addr;
        socklen_t client_len = sizeof(client_addr);
        char first_date[1024];
        char second_date[1024];
        
        while (1) {
            memset(first_date, 0, sizeof(first_date));
            memset(second_date, 0, sizeof(second_date));
            
            int bytes_received = recvfrom(fd, first_date, sizeof(first_date) - 1, 0, 
                                         (struct sockaddr*)&client_addr, &client_len);
            if (bytes_received <= 0) {
                printf("Error receiving data: %d\n", WSAGetLastError());
                continue;
            }
            
            bytes_received = recvfrom(fd, second_date, sizeof(second_date) - 1, 0, 
                                     (struct sockaddr*)&client_addr, &client_len);
            if (bytes_received <= 0) {
                printf("Error receiving data: %d\n", WSAGetLastError());
                continue;
            }

            printf("Received dates: %s and %s\n", first_date, second_date);
            
            char res[1024];
            sprintf(res, "Count of days between two dates is %d", daysBetweenDates(first_date, second_date));
            printf("%s\n", res);
            
            sendto(fd, res, strlen(res) + 1, 0, (struct sockaddr*)&client_addr, client_len);
        }
    }

    closesocket(fd);
    WSACleanup();
    return 0;
}