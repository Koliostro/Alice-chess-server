#include <cstdio>
#include <cstring>
#include <iostream>
#include <sys/socket.h>
#include <sys/types.h>
#include <netinet/in.h>
#include <errno.h>
#include <unistd.h>

int main(int argc, char *argv[]) {

    int port = atoi(argv[1]);

    struct sockaddr_in server_addres;

    bzero((char*)&server_addres, sizeof(server_addres));

    server_addres.sin_family = AF_INET;
    server_addres.sin_addr.s_addr = htonl(INADDR_ANY);
    server_addres.sin_port = htons(port);

    int socket_server = socket(AF_INET, SOCK_STREAM, 0);
    if(socket_server < 0)
    {
        std::cout << strerror(errno) << std::endl;
        exit(0);
    }

    int bind_state = bind(socket_server, (struct sockaddr*) &server_addres, sizeof(server_addres));

    if (bind_state < 0) {
        std::cout << strerror(errno) << std::endl;
        exit(0);
    }

    return 0;
}