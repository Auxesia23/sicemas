<script lang="ts">
    import { auth } from "$lib/states/auth.svelte";
    import { goto } from "$app/navigation";
    import * as v from "valibot";
    import { nipSchema, otpSchema } from "$lib/schemas/auth";

    let nip = $state("");
    let otp = $state("");
    let isOtpSent = $state(false);
    let isLoading = $state(false);
    let errorMessage = $state("");
    let successMessage = $state("");

    const sendOTP = async (e?: Event) => {
        if (e) e.preventDefault();

        const nipValidation = v.safeParse(nipSchema, nip);
        if (!nipValidation.success) {
            errorMessage = nipValidation.issues[0].message;
            return;
        }

        isLoading = true;
        errorMessage = "";

        const err = await auth.login(nip);

        if (err) {
            errorMessage = err;
            isLoading = false;
            return;
        }

        isLoading = false;
        isOtpSent = true;
        successMessage = "Kode OTP telah dikirim ke nomor terdaftar";

        setTimeout(() => {
            successMessage = "";
        }, 3000);
    };

    const handleLogin = async (e?: Event) => {
        if (e) e.preventDefault();

        const otpValidation = v.safeParse(otpSchema, otp);
        if (!otpValidation.success) {
            errorMessage = otpValidation.issues[0].message;
            return;
        }

        isLoading = true;
        errorMessage = "";

        const err = await auth.verifyOtp(otp, nip);

        isLoading = false;

        if (err) {
            errorMessage = err;
            return;
        }

        goto("/dashboard");
    };

    const resetForm = () => {
        nip = "";
        otp = "";
        isOtpSent = false;
        errorMessage = "";
        successMessage = "";
    };
</script>

<svelte:head>
    <title>Login Pengurus - SICEMAS KUA Ciemas</title>
    <meta
        name="description"
        content="Masuk ke dashboard pengurus SICEMAS untuk mengelola data masjid, pesantren, dan majelis taklim di Ciemas."
    />
</svelte:head>

<div
    class="flex min-h-screen items-center justify-center bg-linear-to-b from-base-100 to-base-200 p-4"
>
    <div class="absolute inset-0 overflow-hidden">
        <div
            class="animate-blob absolute -top-40 -right-40 filter h-80 w-80 rounded-full bg-primary opacity-20 mix-blend-multiply blur-xl"
        ></div>
        <div
            class="animate-blob animation-delay-2000 absolute -bottom-40 -left-40 filter h-80 w-80 rounded-full bg-secondary opacity-20 mix-blend-multiply blur-xl"
        ></div>
        <div
            class="animate-blob animation-delay-4000 absolute top-40 left-1/2 filter h-80 w-80 rounded-full bg-accent opacity-20 mix-blend-multiply blur-xl"
        ></div>
    </div>

    <div class="relative z-10 w-full max-w-4xl">
        <div
            class="card overflow-hidden rounded-2xl border border-base-200 shadow-xl"
        >
            <div class="flex flex-col md:flex-row">
                <div
                    class="hidden w-full bg-linear-to-br from-primary to-secondary p-8 md:block md:w-2/5"
                >
                    <div
                        class="relative flex h-full flex-col items-center justify-center overflow-hidden text-center"
                    >
                        <div
                            class="absolute inset-0 flex items-center justify-center"
                        >
                            <div
                                class="h-64 w-64 rounded-full bg-white opacity-10"
                            ></div>
                        </div>
                        <div
                            class="absolute top-10 left-10 h-24 w-24 rotate-45 rounded-lg bg-white opacity-20"
                        ></div>
                        <div
                            class="absolute right-10 bottom-10 h-32 w-32 rounded-full bg-white opacity-15"
                        ></div>
                        <div
                            class="absolute top-1/2 left-1/4 h-16 w-16 rotate-12 rounded-lg bg-white opacity-25"
                        ></div>

                        <div class="relative z-10">
                            <div
                                class="bg-opacity-20 mb-6 inline-block rounded-full bg-white p-4"
                            >
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    class="h-16 w-16 text-white"
                                    fill="none"
                                    viewBox="0 0 24 24"
                                    stroke="currentColor"
                                >
                                    <path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="1.5"
                                        d="M8 14v3m4-3v3m4-3v3M3 21h18M3 10h18M3 7l9-4 9 4M4 10h16v11H4V10z"
                                    />
                                </svg>
                            </div>
                            <h2 class="mb-2 text-2xl font-bold text-white">
                                Sistem Pendataan
                            </h2>
                            <p class="text-opacity-80 text-white">
                                Situs Keagamaan KUA Kecamatan Ciemas
                            </p>
                        </div>
                    </div>
                </div>

                <div class="w-full bg-base-100 md:w-3/5">
                    <div class="card-body flex h-full flex-col p-8">
                        <div class="mb-8 text-center">
                            <div
                                class="mb-4 inline-block rounded-full bg-linear-to-r from-primary to-secondary p-3"
                            >
                                <img
                                    src="/kemenag.webp"
                                    alt="Sicemas KUA Ciemas"
                                    class="h-10 w-10"
                                />
                            </div>
                            <h1
                                class="bg-linear-to-r from-primary to-secondary bg-clip-text text-3xl font-bold text-transparent"
                            >
                                Selamat Datang di Sicemas KUA Ciemas
                            </h1>
                            <p class="mt-2 text-base-content/70">
                                Masuk dengan NIP Anda
                            </p>
                        </div>

                        {#if errorMessage}
                            <div
                                class="mb-4 alert alert-error rounded-xl text-sm"
                            >
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    class="h-6 w-6 shrink-0 stroke-current"
                                    fill="none"
                                    viewBox="0 0 24 24"
                                >
                                    <path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
                                    />
                                </svg>
                                <span>{errorMessage}</span>
                            </div>
                        {/if}

                        {#if successMessage}
                            <div
                                class="mb-4 alert alert-success rounded-xl text-sm"
                            >
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    class="h-6 w-6 shrink-0 stroke-current"
                                    fill="none"
                                    viewBox="0 0 24 24"
                                >
                                    <path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                                    />
                                </svg>
                                <span>{successMessage}</span>
                            </div>
                        {/if}

                        <form
                            class="grow"
                            onsubmit={isOtpSent ? handleLogin : sendOTP}
                        >
                            {#if !isOtpSent}
                                <div class="mb-6">
                                    <label class="label pb-1.5" for="nip">
                                        <span
                                            class="label-text text-xs font-bold uppercase tracking-wider text-base-content/80"
                                            >Nomor Induk Pegawai (NIP/NIK)</span
                                        >
                                    </label>
                                    <input
                                        id="nip"
                                        type="text"
                                        placeholder="Masukkan NIP/NIK Anda"
                                        class="input input-bordered w-full rounded-xl bg-base-200/30 transition-colors focus:bg-base-100"
                                        bind:value={nip}
                                        maxlength="18"
                                        disabled={isLoading}
                                    />
                                </div>
                            {/if}

                            {#if !isOtpSent}
                                <button
                                    type="submit"
                                    class="btn btn-primary mb-6 w-full rounded-xl shadow-lg shadow-primary/20 transition-transform hover:-translate-y-0.5"
                                    disabled={isLoading ||
                                        (nip.length !== 16 &&
                                            nip.length !== 18)}
                                >
                                    {#if isLoading}
                                        <span class="loading loading-spinner"
                                        ></span>
                                        Mengirim OTP...
                                    {:else}
                                        <svg
                                            xmlns="http://www.w3.org/2000/svg"
                                            class="h-5 w-5"
                                            viewBox="0 0 20 20"
                                            fill="currentColor"
                                        >
                                            <path
                                                d="M2.003 5.884L10 9.882l7.997-3.998A2 2 0 0016 4H4a2 2 0 00-1.997 1.884z"
                                            />
                                            <path
                                                d="M18 8.118l-8 4-8-4V14a2 2 0 002 2h12a2 2 0 002-2V8.118z"
                                            />
                                        </svg>
                                        Kirim Kode OTP
                                    {/if}
                                </button>
                            {/if}

                            {#if isOtpSent}
                                <div class="mb-6">
                                    <label class="label pb-1.5" for="otp">
                                        <span
                                            class="label-text text-xs font-bold uppercase tracking-wider text-base-content/80"
                                            >Kode OTP</span
                                        >
                                    </label>
                                    <input
                                        id="otp"
                                        type="text"
                                        bind:value={otp}
                                        maxlength="6"
                                        placeholder="• • • • • •"
                                        class="input input-bordered w-full rounded-xl bg-base-200/30 text-center text-2xl tracking-[0.5em] transition-colors focus:bg-base-100"
                                        disabled={isLoading}
                                        onkeydown={(e) =>
                                            e.key === "Enter" &&
                                            otp.length === 6 &&
                                            handleLogin()}
                                    />
                                    <div class="label pt-2">
                                        <span
                                            class="label-text-alt text-base-content/60"
                                            >Kode telah dikirim ke nomor
                                            terdaftar</span
                                        >
                                    </div>
                                </div>
                            {/if}

                            {#if isOtpSent}
                                <button
                                    type="submit"
                                    class="btn btn-primary mb-4 w-full rounded-xl shadow-lg shadow-primary/20 transition-transform hover:-translate-y-0.5"
                                    disabled={isLoading || otp.length !== 6}
                                >
                                    {#if isLoading}
                                        <span class="loading loading-spinner"
                                        ></span>
                                        Memproses...
                                    {:else}
                                        <svg
                                            xmlns="http://www.w3.org/2000/svg"
                                            class="h-5 w-5"
                                            viewBox="0 0 20 20"
                                            fill="currentColor"
                                        >
                                            <path
                                                fill-rule="evenodd"
                                                d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z"
                                                clip-rule="evenodd"
                                            />
                                        </svg>
                                        Masuk
                                    {/if}
                                </button>
                            {/if}

                            {#if isOtpSent}
                                <div
                                    class="flex items-center justify-between px-1"
                                >
                                    <button
                                        type="button"
                                        class="text-sm font-medium text-base-content/60 hover:text-primary transition-colors"
                                        onclick={resetForm}
                                        disabled={isLoading}
                                    >
                                        Ganti NIP
                                    </button>
                                    <button
                                        type="button"
                                        class="text-sm font-medium text-base-content/60 hover:text-primary transition-colors"
                                        onclick={sendOTP}
                                        disabled={isLoading}
                                    >
                                        Kirim Ulang OTP
                                    </button>
                                </div>
                            {/if}
                        </form>
                    </div>
                </div>
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
