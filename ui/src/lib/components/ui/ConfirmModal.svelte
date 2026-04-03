<script lang="ts">
    interface Props {
        show?: boolean;
        title?: string;
        message?: string;
        confirmText?: string;
        cancelText?: string;
        type?: "danger" | "warning" | "info";
        isProcessing?: boolean;
        onConfirm: () => void;
        onCancel: () => void;
    }

    let {
        show = false,
        title = "Konfirmasi",
        message = "",
        confirmText = "Lanjutkan",
        cancelText = "Batal",
        type = "danger",
        isProcessing = false,
        onConfirm,
        onCancel,
    }: Props = $props();

    let iconColorClass = $derived(
        type === "danger"
            ? "text-error"
            : type === "warning"
              ? "text-warning"
              : "text-info",
    );

    let bgCircleClass = $derived(
        type === "danger"
            ? "bg-error/10"
            : type === "warning"
              ? "bg-warning/10"
              : "bg-info/10",
    );

    let btnConfirmClass = $derived(
        type === "danger"
            ? "btn-error text-white"
            : type === "warning"
              ? "btn-warning text-white"
              : "btn-info text-white",
    );

    let shadowClass = $derived(
        type === "danger"
            ? "shadow-error/20"
            : type === "warning"
              ? "shadow-warning/20"
              : "shadow-info/20",
    );
</script>

<dialog class="modal modal-bottom sm:modal-middle" open={show}>
    <div
        class="modal-box p-0 max-w-md rounded-2xl shadow-2xl overflow-hidden bg-base-100"
    >
        <button
            type="button"
            class="btn absolute top-4 right-4 btn-circle btn-ghost btn-sm"
            onclick={onCancel}
            disabled={isProcessing}
        >
            ✕
        </button>

        <div class="p-6 sm:p-8">
            <div class="flex items-start gap-4 sm:gap-5">
                <div
                    class="flex h-12 w-12 sm:h-14 sm:w-14 shrink-0 items-center justify-center rounded-2xl {bgCircleClass} shadow-sm"
                >
                    {#if type === "danger"}
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="h-6 w-6 sm:h-7 sm:w-7 {iconColorClass}"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                            />
                        </svg>
                    {:else if type === "warning"}
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="h-6 w-6 sm:h-7 sm:w-7 {iconColorClass}"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
                            />
                        </svg>
                    {:else}
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="h-6 w-6 sm:h-7 sm:w-7 {iconColorClass}"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                            />
                        </svg>
                    {/if}
                </div>
                <div class="pt-1 text-left">
                    <h3 class="text-xl font-extrabold text-base-content">
                        {title}
                    </h3>
                    <p
                        class="mt-2 text-sm font-medium text-base-content/60 leading-relaxed"
                    >
                        {message}
                    </p>
                </div>
            </div>
        </div>

        <div
            class="bg-base-200/40 px-6 py-5 border-t border-base-200 sm:px-8 flex flex-col-reverse sm:flex-row justify-end gap-3"
        >
            <button
                type="button"
                class="btn btn-ghost rounded-xl hover:bg-base-300 w-full sm:w-auto"
                onclick={onCancel}
                disabled={isProcessing}
            >
                {cancelText}
            </button>
            <button
                type="button"
                class="btn {btnConfirmClass} rounded-xl px-6 shadow-lg {shadowClass} transition-transform hover:-translate-y-0.5 w-full sm:w-auto"
                onclick={onConfirm}
                disabled={isProcessing}
            >
                {#if isProcessing}
                    <span class="loading loading-sm loading-spinner"></span> Memproses...
                {:else}
                    {confirmText}
                {/if}
            </button>
        </div>
    </div>

    <button
        type="button"
        class="modal-backdrop bg-base-content/20 backdrop-blur-sm cursor-default"
        onclick={onCancel}
        aria-label="Tutup modal"
        disabled={isProcessing}
    ></button>
</dialog>
