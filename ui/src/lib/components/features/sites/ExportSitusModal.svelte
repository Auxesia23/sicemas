<script lang="ts">
    import type { JenisSitus } from "$lib/schemas/sites";

    interface Props {
        show: boolean;
        jenisSitus: JenisSitus[];
        isExporting: boolean;
        onClose: () => void;
        onExport: (selectedType: string) => void;
    }

    let { show, jenisSitus, isExporting, onClose, onExport }: Props = $props();

    let exportSelectedType = $state("");

    $effect(() => {
        if (show) {
            exportSelectedType = "";
        }
    });
</script>

<dialog class="modal modal-bottom sm:modal-middle" open={show}>
    <div
        class="modal-box p-0 max-w-lg rounded-2xl shadow-2xl overflow-hidden bg-base-100"
    >
        <div
            class="bg-base-200/40 px-6 py-5 border-b border-base-200 sm:px-8 flex justify-between items-center"
        >
            <div class="flex items-center gap-4">
                <div
                    class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-linear-to-br from-success/20 to-primary/20 text-success shadow-sm"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="h-5 w-5"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"
                        />
                    </svg>
                </div>
                <h3 class="text-xl font-extrabold text-base-content">
                    Export Data ke Excel
                </h3>
            </div>
            <button
                type="button"
                class="btn btn-circle btn-ghost btn-sm"
                onclick={onClose}
                disabled={isExporting}
            >
                ✕
            </button>
        </div>

        <div class="px-6 py-6 sm:px-8 space-y-5">
            <p class="text-sm font-medium text-base-content/60 leading-relaxed">
                Pilih jenis situs yang ingin diekspor. Data diekspor per jenis
                agar format kolom pada file Excel (seperti detail fasilitas dan
                pengurus) tetap rapi dan terstruktur.
            </p>

            <div class="form-control w-full">
                <label class="label pb-1.5 pt-0" for="exportType">
                    <span class="label-text font-bold text-base-content/80">
                        Jenis Situs <span class="text-error">*</span>
                    </span>
                </label>
                <select
                    id="exportType"
                    class="select select-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors focus:border-success"
                    bind:value={exportSelectedType}
                    disabled={isExporting}
                >
                    <option value="" disabled selected={!exportSelectedType}
                        >Pilih Jenis Situs...</option
                    >
                    {#each jenisSitus as jenis (jenis.id)}
                        <option value={jenis.nama}>{jenis.nama}</option>
                    {/each}
                </select>
            </div>
        </div>

        <div
            class="bg-base-200/40 px-6 py-5 border-t border-base-200 sm:px-8 flex flex-col-reverse sm:flex-row justify-end gap-3"
        >
            <button
                type="button"
                class="btn btn-ghost rounded-xl hover:bg-base-300 w-full sm:w-auto"
                onclick={onClose}
                disabled={isExporting}
            >
                Batal
            </button>
            <button
                type="button"
                class="btn btn-success rounded-xl px-8 shadow-lg shadow-success/20 transition-transform hover:-translate-y-0.5 w-full sm:w-auto text-white"
                onclick={() => onExport(exportSelectedType)}
                disabled={!exportSelectedType || isExporting}
            >
                {#if isExporting}
                    <span class="loading loading-sm loading-spinner"></span> Menyiapkan...
                {:else}
                    Download Excel
                {/if}
            </button>
        </div>
    </div>

    <button
        type="button"
        class="modal-backdrop bg-base-content/20 backdrop-blur-sm"
        aria-label="Close"
        onclick={onClose}
        disabled={isExporting}
    ></button>
</dialog>
