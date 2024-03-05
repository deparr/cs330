#include <bits/types/sigset_t.h>
#include <signal.h>
#include <stdio.h>
#include <threads.h>
#include <stdlib.h>
#include <unistd.h>

void catch_sigsev(int signum) {
	sigset_t mask;
	sigemptyset(&mask);
	sigaddset(&mask, SIGSEGV);
	sigprocmask(SIG_SETMASK, &mask, NULL);
	puts("Caught SIGSEGV!");
}

void catch_sigint(int signum) {
	puts("caught sigint");
}


int main() {
	struct sigaction sigact;
	sigact.sa_flags = SA_RESTART;
	sigact.sa_handler = catch_sigsev;
	sigaction(SIGSEGV, &sigact, NULL);
	sigact.sa_handler = catch_sigint;
	sigaction(SIGINT, &sigact, NULL);

	// sleep(10);
	int* b = NULL;
	int a = *(int*)b;
	puts("After sigsev!");
}
