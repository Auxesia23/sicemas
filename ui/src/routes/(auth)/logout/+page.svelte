<script lang="ts">
    import { onMount, onDestroy } from "svelte";
    import { auth } from "$lib/states/auth.svelte";
    import { goto } from "$app/navigation";

    let countdown = $state(5);
    let timer: ReturnType<typeof setInterval>;

    const logoutNow = async () => {
        if (timer) {
            clearInterval(timer);
        }
        await auth.logout();
    };

    onMount(() => {
        timer = setInterval(() => {
            countdown--;
            if (countdown <= 0) {
                clearInterval(timer);
                logoutNow();
            }
        }, 1000);
    });

    onDestroy(() => {
        if (timer) {
            clearInterval(timer);
        }
    });
</script>

<svelte:head>
    <title>Keluar - SICEMAS</title>
</svelte:head>

<div
    class="relative flex min-h-screen items-center justify-center bg-linear-to-b from-base-100 to-base-200 p-4 overflow-hidden"
>
    <div
        class="pointer-events-none absolute inset-0 flex items-center justify-center z-0"
    >
        <div
            class="animate-blob filter h-64 w-64 rounded-full bg-primary opacity-20 mix-blend-multiply blur-xl"
        ></div>
        <div
            class="animate-blob animation-delay-2000 filter ml-32 h-64 w-64 rounded-full bg-error opacity-20 mix-blend-multiply blur-xl"
        ></div>
        <div
            class="animate-blob animation-delay-4000 filter -mt-32 h-64 w-64 rounded-full bg-warning opacity-20 mix-blend-multiply blur-xl"
        ></div>
    </div>

    <div
        class="card w-full max-w-md border border-base-200/50 bg-base-100/80 backdrop-blur-md shadow-2xl z-10"
    >
        <div class="card-body items-center text-center p-8">
            <div
                class="mb-6 flex h-20 w-20 items-center justify-center rounded-full bg-error/10 shadow-inner"
            >
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-10 w-10 text-error"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
                    />
                </svg>
            </div>

            <h2 class="text-3xl font-bold text-gray-800 mb-2">
                Mengakhiri Sesi
            </h2>
            <p class="mb-8 text-base-content/70">
                Sesi Anda pada aplikasi SICEMAS akan segera diakhiri. Apakah
                Anda yakin ingin keluar sekarang?
            </p>

            <div
                class="mb-8 w-full rounded-xl bg-base-200/50 p-4 border border-base-300"
            >
                <p class="text-sm font-medium text-base-content/70">
                    Otomatis keluar dalam
                </p>
                <div class="mt-1 flex items-baseline justify-center gap-1">
                    <span
                        class="text-4xl font-black text-error font-mono tabular-nums leading-none"
                    >
                        {countdown}
                    </span>
                    <span class="text-sm font-bold text-base-content/50"
                        >detik</span
                    >
                </div>
            </div>

            <div class="w-full flex flex-col gap-3">
                <button
                    type="button"
                    class="btn w-full btn-error text-white shadow-lg hover:shadow-error/30 hover:-translate-y-0.5 transition-all"
                    onclick={logoutNow}
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="mr-2 h-5 w-5"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
                        />
                    </svg>
                    Ya, Keluar Sekarang
                </button>

                <button
                    type="button"
                    class="btn w-full btn-outline bg-base-100 hover:bg-base-200 hover:text-base-content"
                    onclick={() => {
                        if (timer) clearInterval(timer);
                        goto("/dashboard");
                    }}
                >
                    Batal, Kembali ke Dashboard
                </button>
            </div>
        </div>
    </div>
</div>

<style>
    @keyframes blob {
        0% {
            transform: translate(0px, 0px) scale(1);
        }
        33% {
            transform: translate(30px, -50px) scale(1.1);
        }
        66% {
            transform: translate(-20px, 20px) scale(0.9);
        }
        100% {
            transform: translate(0px, 0px) scale(1);
        }
    }
    .animate-blob {
        animation: blob 7s infinite;
    }
    .animation-delay-2000 {
        animation-delay: 2s;
    }
    .animation-delay-4000 {
        animation-delay: 4s;
    }
</style>
