<script lang="ts">
    import { fly } from "svelte/transition";

    interface Props {
        show?: boolean;
        message?: string;
        type?: "success" | "error" | "warning" | "info";
        onclose?: () => void;
    }

    let {
        show = false,
        message = "",
        type = "success",
        onclose,
    }: Props = $props();

    let config = $derived.by(() => {
        switch (type) {
            case "success":
                return {
                    title: "Berhasil!",
                    iconClass: "bg-success/10 text-success",
                    shadowClass: "shadow-success/20",
                    icon: "M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z",
                };
            case "error":
                return {
                    title: "Gagal!",
                    iconClass: "bg-error/10 text-error",
                    shadowClass: "shadow-error/20",
                    icon: "M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z",
                };
            case "warning":
                return {
                    title: "Peringatan!",
                    iconClass: "bg-warning/10 text-warning",
                    shadowClass: "shadow-warning/20",
                    icon: "M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z",
                };
            case "info":
            default:
                return {
                    title: "Informasi",
                    iconClass: "bg-info/10 text-info",
                    shadowClass: "shadow-info/20",
                    icon: "M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z",
                };
        }
    });
</script>

{#if show}
    <div class="toast toast-top toast-end z-9999 pt-20 sm:pt-6 overflow-hidden">
        <div
            class="flex items-start gap-4 p-4 min-w-75 max-w-sm rounded-2xl bg-base-100/90 backdrop-blur-xl border border-base-200 shadow-2xl {config.shadowClass}"
            transition:fly={{ x: 100, duration: 400, opacity: 0 }}
        >
            <div
                class="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl {config.iconClass}"
            >
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-6 w-6"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d={config.icon}
                    />
                </svg>
            </div>

            <div class="flex-1 pt-1 min-w-0">
                <h3
                    class="text-base font-extrabold text-base-content tracking-tight"
                >
                    {config.title}
                </h3>
                <p
                    class="mt-1 text-sm font-medium text-base-content/60 leading-relaxed"
                >
                    {message}
                </p>
            </div>

            {#if onclose}
                <button
                    type="button"
                    class="btn btn-ghost btn-sm btn-circle h-8 w-8 min-h-0 text-base-content/40 hover:bg-base-200 hover:text-base-content transition-colors"
                    onclick={onclose}
                    aria-label="Tutup"
                >
                    ✕
                </button>
            {/if}
        </div>
    </div>
{/if}
