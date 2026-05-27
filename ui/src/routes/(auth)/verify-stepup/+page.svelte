<script lang="ts">
    import { apiFetch } from "$lib/api";
    import { goto } from "$app/navigation";
    import { page } from "$app/stores";
    import { auth } from "$lib/states/auth.svelte";

    let otp = $state("");
    let loading = $state(false);
    let resending = $state(false);
    let errorMessage = $state("");
    let successMessage = $state("");

    async function verifyStepUp() {
        loading = true;
        errorMessage = "";
        successMessage = "";

        try {
            const res = await apiFetch("/api/auth/verify-stepup", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ otp }),
            });

            if (res.ok) {
                auth.isStepUpRequired = false;
                const nextUrl =
                    $page.url.searchParams.get("next") || "/dashboard";
                goto(nextUrl);
            } else {
                const errText = await res.json();
                errorMessage =
                    errText.error || "OTP tidak valid atau kadaluarsa.";
            }
        } catch (error) {
            errorMessage = "Terjadi kesalahan jaringan.";
        } finally {
            loading = false;
        }
    }

    async function resendOtp() {
        resending = true;
        errorMessage = "";
        successMessage = "";

        try {
            const res = await apiFetch("/api/auth/resend-stepup", {
                method: "POST",
            });
            if (res.ok) {
                successMessage = "OTP baru telah dikirim ke WhatsApp Anda.";
            } else {
                errorMessage = "Gagal mengirim ulang OTP.";
            }
        } catch {
            errorMessage = "Terjadi kesalahan jaringan.";
        } finally {
            resending = false;
        }
    }

    const cancel = async () => {
        auth.isStepUpRequired = false;
        await auth.logout();
    };
</script>

<div class="flex min-h-[85vh] items-center justify-center p-4">
    <div
        class="card w-full max-w-lg rounded-2xl border border-base-200 bg-base-100 shadow-xl shadow-base-200/40"
    >
        <div class="card-body p-8">
            <div class="mb-6">
                <div class="flex items-center gap-3">
                    <div
                        class="h-8 w-1.5 shrink-0 rounded-full bg-linear-to-b from-primary to-secondary"
                    ></div>
                    <h1
                        class="text-3xl font-extrabold tracking-tight text-base-content"
                    >
                        Verifikasi Sesi
                    </h1>
                </div>
                <p class="mt-1.5 pl-4 text-sm font-medium text-base-content/60">
                    Zero Trust Architecture - MFA
                </p>
            </div>

            <div class="mb-6 flex items-start gap-4">
                <div
                    class="flex h-12 w-12 shrink-0 items-center justify-center rounded-2xl bg-linear-to-br from-primary/20 to-secondary/20 text-primary shadow-sm"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="24"
                        height="24"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2.5"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                    >
                        <path
                            d="M20 13c0 5-3.5 7.5-7.66 8.95a1 1 0 0 1-.67-.01C7.5 20.5 4 18 4 13V6a1 1 0 0 1 1-1c2-1 4-2 7-2 2.5 0 4.5 1.2 7 2a1 1 0 0 1 1 1v7z"
                        />
                        <path d="M9 12l2 2 4-4" />
                    </svg>
                </div>
                <div class="pt-1">
                    <p class="text-sm font-medium text-base-content/70">
                        Sistem mendeteksi aktivitas dari lokasi atau perangkat
                        baru. Demi keamanan data KUA, silakan masukkan OTP dari
                        WhatsApp Anda.
                    </p>
                </div>
            </div>

            {#if errorMessage}
                <div
                    class="mb-5 flex items-center rounded-xl border border-error/20 bg-error/10 px-4 py-3 text-sm font-medium text-error"
                >
                    <div
                        class="mr-2 h-1.5 w-1.5 shrink-0 rounded-full bg-error"
                    ></div>
                    {errorMessage}
                </div>
            {/if}

            {#if successMessage}
                <div
                    class="mb-5 flex items-center rounded-xl border border-success/20 bg-success/10 px-4 py-3 text-sm font-medium text-success"
                >
                    <div
                        class="mr-2 h-1.5 w-1.5 shrink-0 rounded-full bg-success"
                    ></div>
                    {successMessage}
                </div>
            {/if}

            <div class="form-control w-full">
                <label class="label pb-1.5" for="otp-input">
                    <span
                        class="label-text text-xs font-bold uppercase tracking-wider text-base-content/80"
                        >Kode OTP</span
                    >
                </label>
                <input
                    id="otp-input"
                    type="text"
                    bind:value={otp}
                    maxlength="6"
                    placeholder="• • • • • •"
                    class="input input-bordered w-full rounded-xl bg-base-200/30 text-center text-2xl tracking-[0.5em] transition-colors focus:bg-base-100"
                    disabled={loading}
                    onkeydown={(e) =>
                        e.key === "Enter" && otp.length === 6 && verifyStepUp()}
                />
            </div>

            <div class="mt-4 flex justify-center">
                <button
                    onclick={resendOtp}
                    disabled={resending || loading}
                    class="flex items-center gap-1.5 text-xs font-bold text-base-content/50 transition-colors hover:text-primary disabled:opacity-50"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="14"
                        height="14"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        class={resending ? "animate-spin" : ""}
                    >
                        <path
                            d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"
                        />
                        <path d="M3 3v5h5" />
                    </svg>
                    {resending ? "Mengirim ulang..." : "Kirim ulang OTP"}
                </button>
            </div>

            <div class="card-actions mt-8 flex items-center justify-end gap-3">
                <button
                    onclick={cancel}
                    class="btn btn-ghost rounded-xl text-base-content/70 hover:bg-primary/10 hover:text-primary"
                    disabled={loading}
                >
                    Kembali ke Login
                </button>
                <button
                    onclick={verifyStepUp}
                    disabled={loading || otp.length < 6}
                    class="btn btn-primary rounded-xl px-8 shadow-lg shadow-primary/20 transition-transform hover:-translate-y-0.5 disabled:opacity-50"
                >
                    {#if loading}
                        <span class="loading loading-spinner loading-sm"></span>
                        Memverifikasi...
                    {:else}
                        Lanjutkan
                    {/if}
                </button>
            </div>
        </div>
    </div>
</div>
