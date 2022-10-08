#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>
#include <fcntl.h>

int main(int argc, char *argv[]) {
    int fd = open("./c/check.txt", O_CREAT | O_RDWR | O_TRUNC, S_IRWXU);
    int rc = fork();
    if (rc < 0) {
        fprintf(stderr, "fork failed\n");
    } else if (rc == 0) {
        char *buf = strdup("child\n");
        int error = write(fd, buf, sizeof(char) * strlen(buf));
        printf("child error: %d\n", error == -1 ? 1 : 0);
    } else {
        char *buf = strdup("parent\n");
        int error = write(fd, buf, sizeof(char) * strlen(buf));
        printf("parent error: %d\n", error == -1 ? 1 : 0);
        wait(NULL);
    }
    return 0;
}