#include <stdio.h> 
#include <string.h>

int main(int argc, char** argv) {
    printf("Content-Type: text/plain\n\n");
    char* subj;
    if ( argc > 1 && strlen(argv[1]) != 0 ) {
        subj = argv [1];
    } else {
        subj = "world";
    }
    printf("Hello %s\n", subj);
    return 0;
}
