#include <winsock2.h>
#include <stdio.h>
#include <string.h>
#include <ws2tcpip.h>

int main(int argc, char** argv) {
    int use_udp = 0;
    char server_ip[256] = "127.0.0.1"; // Default localhost
    
    for (int i = 1; i < argc; i++) {
        if (strcmp(argv[i], "-u") == 0) {
            use_udp = 1;
        }
        else if (strcmp(argv[i], "-ip") == 0 && i + 1 < argc) {
            strcpy(server_ip, argv[i + 1]);
            i++;
        }
        else if (strcmp(argv[i], "-h") == 0 || strcmp(argv[i], "--help") == 0) {
            printf("Usage: %s [options]\n", argv[0]);
            printf("Options:\n");
            printf("  -u              Use UDP protocol (default: TCP)\n");
            printf("  -ip <address>   Server IP address (default: 127.0.0.1)\n");
            printf("  -h, --help      Show this help message\n");
            WSACleanup();
            return 0;
        }
    }

    // Display configuration
    printf("Protocol: %s\n", use_udp ? "UDP" : "TCP");
    printf("Server IP: %s\n", server_ip);
    printf("Server Port: 8080\n");

    WSADATA wsa;
    if (WSAStartup(MAKEWORD(2, 2), &wsa) != 0) {
        printf("error: WSA failed\n");
        return 1;
    }

    int sock;
    if (use_udp) {
        sock = socket(AF_INET, SOCK_DGRAM, 0);
    } else {
        sock = socket(AF_INET, SOCK_STREAM, 0);
    }

    if (sock == INVALID_SOCKET) {
        printf("error: Socket creation failed: %d\n", WSAGetLastError());
        WSACleanup();
        return 1;
    }

    struct sockaddr_in server_addr;
    server_addr.sin_family = AF_INET;
    server_addr.sin_port = htons(8080);
    
    if (inet_pton(AF_INET, server_ip, &server_addr.sin_addr) <= 0) {
        printf("error: Invalid address: %s\n", server_ip);
        closesocket(sock);
        WSACleanup();
        return 1;
    }

    if (!use_udp) {
        printf("Connecting to %s:8080...\n", server_ip);
        if (connect(sock, (struct sockaddr*)&server_addr, sizeof(server_addr)) == SOCKET_ERROR) {
            printf("error: Connect failed: %d\n", WSAGetLastError());
            closesocket(sock);
            WSACleanup();
            return 1;
        }
        printf("Connected to TCP server\n");
    } else {
        printf("UDP client ready to send to %s:8080\n", server_ip);
    }

    char first_date[1024];
    char second_date[1024];
    char response[1024];

    while (1) {
        // Get first date from user
        printf("\nEnter first date (DD.MM.YYYY) or 'quit' to exit: ");
        fgets(first_date, sizeof(first_date), stdin);
        first_date[strcspn(first_date, "\n")] = 0; // Remove newline
        
        if (strcmp(first_date, "quit") == 0) {
            break;
        }

        // Get second date from user
        printf("Enter second date (DD.MM.YYYY): ");
        fgets(second_date, sizeof(second_date), stdin);
        second_date[strcspn(second_date, "\n")] = 0; // Remove newline

        if (use_udp) {
            // UDP: Send dates and receive response
            socklen_t addr_len = sizeof(server_addr);
            
            // Send first date
            if (sendto(sock, first_date, strlen(first_date) + 1, 0, 
                      (struct sockaddr*)&server_addr, sizeof(server_addr)) == SOCKET_ERROR) {
                printf("error: Send failed: %d\n", WSAGetLastError());
                continue;
            }
            printf("Sent first date: %s\n", first_date);
            
            // Small delay to ensure packets are sent separately
            Sleep(100);
            
            // Send second date
            if (sendto(sock, second_date, strlen(second_date) + 1, 0, 
                      (struct sockaddr*)&server_addr, sizeof(server_addr)) == SOCKET_ERROR) {
                printf("error: Send failed: %d\n", WSAGetLastError());
                continue;
            }
            printf("Sent second date: %s\n", second_date);
            
            // Receive response
            memset(response, 0, sizeof(response));
            int bytes_received = recvfrom(sock, response, sizeof(response) - 1, 0,
                                         (struct sockaddr*)&server_addr, &addr_len);
            if (bytes_received <= 0) {
                printf("error: Receive failed: %d\n", WSAGetLastError());
                continue;
            }
            
            printf("Server response: %s\n", response);
            
        } else {
            // TCP: Send dates and receive response
            // Send first date
            if (send(sock, first_date, strlen(first_date) + 1, 0) == SOCKET_ERROR) {
                printf("error: Send failed: %d\n", WSAGetLastError());
                break;
            }
            printf("Sent first date: %s\n", first_date);
            
            // Send second date
            if (send(sock, second_date, strlen(second_date) + 1, 0) == SOCKET_ERROR) {
                printf("error: Send failed: %d\n", WSAGetLastError());
                break;
            }
            printf("Sent second date: %s\n", second_date);
            
            // Receive response
            memset(response, 0, sizeof(response));
            int bytes_received = recv(sock, response, sizeof(response) - 1, 0);
            if (bytes_received <= 0) {
                printf("Server disconnected or error occurred\n");
                break;
            }
            
            printf("Server response: %s\n", response);
        }
    }

    printf("Closing connection...\n");
    closesocket(sock);
    WSACleanup();
    return 0;
}