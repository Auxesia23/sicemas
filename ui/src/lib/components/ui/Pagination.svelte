<script lang="ts">
    interface Props {
        currentPage: number;
        itemsPerPage: number;
        totalItems: number;
        onPageChange?: (page: number) => void;
    }

    let {
        currentPage = $bindable(1),
        itemsPerPage = $bindable(10),
        totalItems,
        onPageChange,
    }: Props = $props();

    let totalPages = $derived(Math.ceil(totalItems / itemsPerPage) || 1);

    function goToPage(page: number) {
        if (page >= 1 && page <= totalPages) {
            currentPage = page;
            if (onPageChange) onPageChange(page);
        }
    }
</script>

<div
    class="mt-6 mb-2 flex flex-col items-center justify-between gap-4 px-4 sm:flex-row"
>
    <div class="flex items-center gap-2" title="Pilih jumlah baris per halaman">
        <div
            class="flex h-8 w-8 items-center justify-center rounded-md bg-base-200 text-base-content/70"
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
                    d="M4 6h16M4 12h16M4 18h7"
                />
            </svg>
        </div>
        <select
            class="select-bordered select select-sm font-medium"
            bind:value={itemsPerPage}
        >
            <option value={5}>5</option>
            <option value={10}>10</option>
            <option value={25}>25</option>
            <option value={50}>50</option>
        </select>
    </div>

    <div class="join">
        <button
            type="button"
            class="btn join-item px-2 btn-sm"
            disabled={currentPage === 1}
            onclick={() => goToPage(currentPage - 1)}
            title="Halaman Sebelumnya"
        >
            <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-5 w-5"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                stroke-width="2"
            >
                <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M15 19l-7-7 7-7"
                />
            </svg>
        </button>

        <button
            type="button"
            class="no-animation btn pointer-events-none join-item flex items-center justify-center bg-base-200 px-4 font-medium text-base-content/80 btn-sm"
        >
            Halaman {currentPage} / {totalPages}
        </button>

        <button
            type="button"
            class="btn join-item px-2 btn-sm"
            disabled={currentPage === totalPages}
            onclick={() => goToPage(currentPage + 1)}
            title="Halaman Selanjutnya"
        >
            <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-5 w-5"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                stroke-width="2"
            >
                <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M9 5l7 7-7 7"
                />
            </svg>
        </button>
    </div>
</div>
