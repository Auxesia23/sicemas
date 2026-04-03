<script lang="ts">
    import { page } from "$app/state";

    let status = $derived(page.status);
    let errorMessage = $derived(
        page.error?.message ||
            "Terjadi kesalahan yang tidak diketahui pada sistem.",
    );

    let title = $derived.by(() => {
        switch (status) {
            case 404:
                return "Halaman Tidak Ditemukan";
            case 500:
                return "Server Mengalami Gangguan";
            case 403:
                return "Akses Dibatasi";
            case 401:
                return "Sesi Habis / Belum Login";
            default:
                return "Terjadi Kesalahan Sistem";
        }
    });
</script>

<div
    class="relative flex min-h-screen items-center justify-center overflow-hidden bg-base-200 p-4 z-0"
>
    <div
        class="absolute -top-20 -left-20 h-72 w-72 rounded-full bg-primary/10 blur-3xl"
    ></div>
    <div
        class="absolute -bottom-20 -right-20 h-72 w-72 rounded-full bg-secondary/10 blur-3xl"
    ></div>

    <div
        class="card w-full max-w-lg rounded-2xl border border-base-200 bg-base-100/80 backdrop-blur-xl shadow-xl shadow-base-200/40 relative z-10 transition-all duration-300 hover:shadow-2xl hover:shadow-base-200/50"
    >
        <div class="card-body items-center text-center p-8 sm:p-10">
            <div
                class="mb-6 flex h-32 w-32 items-center justify-center rounded-full bg-linear-to-br from-primary/20 to-secondary/20 shadow-inner"
            >
                <h1
                    class="text-6xl font-black tracking-tighter text-primary drop-shadow-sm"
                >
                    {status}
                </h1>
            </div>

            <div class="flex items-center justify-center gap-3">
                <div
                    class="h-8 w-1.5 rounded-full bg-linear-to-b from-primary to-secondary"
                ></div>
                <h2
                    class="text-3xl font-extrabold tracking-tight text-base-content"
                >
                    {title}
                </h2>
            </div>

            <p class="mt-4 mb-8 text-sm font-medium text-base-content/60">
                {errorMessage}
            </p>

            <div class="flex w-full flex-col justify-center gap-3 sm:flex-row">
                {#if status === 401}
                    <a
                        href="/login"
                        class="btn btn-primary rounded-xl px-8 shadow-lg shadow-primary/20 transition-transform hover:-translate-y-0.5 w-full sm:w-auto"
                    >
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="mr-1.5 h-5 w-5"
                            viewBox="0 0 20 20"
                            fill="currentColor"
                        >
                            <path
                                fill-rule="evenodd"
                                d="M3 3a1 1 0 011 1v12a1 1 0 11-2 0V4a1 1 0 011-1zm7.707 3.293a1 1 0 010 1.414L6.414 9H17a1 1 0 110 2H6.414l4.293 4.293a1 1 0 01-1.414 1.414l-6-6a1 1 0 010-1.414l6-6a1 1 0 011.414 0z"
                                clip-rule="evenodd"
                            />
                        </svg>
                        Ke Halaman Login
                    </a>
                {:else if status === 500}
                    <button
                        class="btn btn-ghost border border-base-200 bg-base-100/50 rounded-xl px-8 text-base-content/80 transition-transform hover:-translate-y-0.5 hover:bg-base-200 hover:text-base-content w-full sm:w-auto"
                        onclick={() => window.location.reload()}
                    >
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="mr-1.5 h-5 w-5"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
                            />
                        </svg>
                        Coba Lagi
                    </button>
                    <a
                        href="/dashboard"
                        class="btn btn-primary rounded-xl px-8 shadow-lg shadow-primary/20 transition-transform hover:-translate-y-0.5 w-full sm:w-auto"
                    >
                        Ke Dashboard
                    </a>
                {:else}
                    <button
                        class="btn btn-ghost border border-base-200 bg-base-100/50 rounded-xl px-8 text-base-content/80 transition-transform hover:-translate-y-0.5 hover:bg-base-200 hover:text-base-content w-full sm:w-auto"
                        onclick={() => window.history.back()}
                    >
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="mr-1.5 h-5 w-5"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M10 19l-7-7m0 0l7-7m-7 7h18"
                            />
                        </svg>
                        Kembali
                    </button>
                    <a
                        href="/dashboard"
                        class="btn btn-primary rounded-xl px-8 shadow-lg shadow-primary/20 transition-transform hover:-translate-y-0.5 w-full sm:w-auto"
                    >
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="mr-1.5 h-5 w-5"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"
                            />
                        </svg>
                        Ke Dashboard
                    </a>
                {/if}
            </div>
        </div>
    </div>
</div>
