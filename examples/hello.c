#include <stdio.h> 
#include <string.h>

#define KEY_NAME "name="
#define KEY_NAME_LEN ((sizeof(KEY_NAME))-1)

int main(int argc, char** argv) {
    printf("Content-Type: text/plain\n\n");
    char* subj = "world";

    for (int i = 1; i < argc; i++) {
        // look for for name=<value> in the list of <key>=<value> pairs
        if (memcmp(KEY_NAME, argv[i], KEY_NAME_LEN) == 0) {
            subj = argv[i] + KEY_NAME_LEN;
        }
    }

    printf("Hello %s!\n", subj);
    return 0;
}
