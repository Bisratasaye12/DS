// #include <sys/types.h>
// #include <unistd.h>
// #include <stdio.h>




// int main() {
//     pid_t pid;
//     pid = fork();
//     if (pid < 0) {
//         printf("Fork failed\n");
//         return 1;
//     }
//     if (pid == 0) {
//         printf("Child process\n");
//     } else {
//         printf("Parent process\n");
//     }

//     printf("pid %d",getpid());
//     printf("pid %d",getppid());
//     printf("pid %d",getpid());
//     printf("pid %d",getppid());
//     return 0;
// }



// pipe

#include <stdio.h>
#include <unistd.h>
#include <string.h>

int main() {
    int fd[2];  // File descriptors for the pipe
    pid_t pid;
    char write_msg[] = "Hello from parent!";
    char read_msg[100];

    // Create the pipe
    if (pipe(fd) == -1) {
        perror("pipe");
        return 1;
    }

    // Fork a child process
    pid = fork();

    if (pid < 0) {
        perror("fork");
        return 1;
    }

    if (pid == 0) {  // Child process
        close(fd[1]);  // Close unused write end

        // Read message from the pipe
        read(fd[0], read_msg, sizeof(read_msg));
        printf("Child received: %s\n", read_msg);

        close(fd[0]);  // Close read end
    } else {  // Parent process
        close(fd[0]);  // Close unused read end

        // Write message to the pipe
        write(fd[1], write_msg, strlen(write_msg) + 1);

        close(fd[1]);  // Close write end
    }

    return 0;
}
