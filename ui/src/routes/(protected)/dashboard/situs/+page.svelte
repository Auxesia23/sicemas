<script lang="ts">
    import type { SitusKeagamaan } from "$lib/schemas/sites";
    import type { PageData } from "./$types";
    import { goto, invalidate } from "$app/navigation";
    import { apiFetch } from "$lib/api";
    import { auth } from "$lib/states/auth.svelte";

    import ConfirmModal from "$lib/components/ui/ConfirmModal.svelte";
    import Toast from "$lib/components/ui/Toast.svelte";
    import SitusTableActions from "$lib/components/features/sites/SitusTableActions.svelte";
    import Pagination from "$lib/components/ui/Pagination.svelte";
    import ExportSitusModal from "$lib/components/features/sites/ExportSitusModal.svelte";

    let { data }: { data: PageData } = $props();

    let { situs, jenisSitus } = $derived(data);

    let isRefreshing = $state(false);

    let searchTerm = $state("");
    let statusFilter = $state("all");
    let typeFilter = $state("all");

    let currentPage = $state(1);
    let itemsPerPage = $state(10);

    let showDeleteModal = $state(false);
    let siteToDelete = $state<SitusKeagamaan | null>(null);
    let isDeleting = $state(false);

    let showToast = $state(false);
    let toastMessage = $state("");
    let toastType = $state<"success" | "error" | "info" | "warning">("success");

    let showExportModal = $state(false);
    let isExporting = $state(false);

    let filteredSites = $derived(
        situs.filter((site) => {
            const searchLower = searchTerm.toLowerCase();

            const siteIdKemenag = site.situs_id
                ? String(site.situs_id).toLowerCase()
                : "";
            const siteName = site.nama ? site.nama.toLowerCase() : "";
            const siteLocation = site.lokasi ? site.lokasi.toLowerCase() : "";
            const sitePendata = site.pendata ? site.pendata.toLowerCase() : "";

            const matchesSearch =
                siteName.includes(searchLower) ||
                siteLocation.includes(searchLower) ||
                sitePendata.includes(searchLower) ||
                siteIdKemenag.includes(searchLower);

            const matchesStatus =
                statusFilter === "all" || site.status === statusFilter;
            const matchesType =
                typeFilter === "all" || site.jenis === typeFilter;

            return matchesSearch && matchesStatus && matchesType;
        }),
    );

    let paginatedSites = $derived(
        filteredSites.slice(
            (currentPage - 1) * itemsPerPage,
            currentPage * itemsPerPage,
        ),
    );

    $effect(() => {
        searchTerm;
        statusFilter;
        typeFilter;
        itemsPerPage;
        currentPage = 1;
    });

    function formatDate(dateString: string | null): string {
        if (!dateString) return "-";
        const date = new Date(dateString);
        return date.toLocaleDateString("id-ID", {
            day: "numeric",
            month: "short",
            year: "numeric",
            hour: "2-digit",
            minute: "2-digit",
        });
    }

    function getStatusBadgeClass(status: string): string {
        switch (status) {
            case "menunggu":
                return "badge-warning";
            case "terverifikasi":
                return "badge-success";
            case "ditolak":
                return "badge-error";
            default:
                return "badge-ghost";
        }
    }

    async function handleRefresh() {
        isRefreshing = true;
        await invalidate("data:situs");
        isRefreshing = false;
    }

    function confirmDelete(site: SitusKeagamaan) {
        siteToDelete = site;
        showDeleteModal = true;
    }

    async function handleDelete() {
        if (!siteToDelete) return;

        isDeleting = true;
        try {
            const res = await apiFetch("/api/situs/" + siteToDelete.id, {
                method: "DELETE",
            });

            if (res.ok) {
                showDeleteModal = false;
                toastMessage = "Data situs berhasil dihapus";
                toastType = "success";
                showToast = true;

                await invalidate("data:situs");

                setTimeout(() => {
                    showToast = false;
                }, 3000);
            } else {
                const errText = await res.text();
                toastMessage = "Gagal menghapus data: " + errText;
                toastType = "error";
                showToast = true;
                setTimeout(() => {
                    showToast = false;
                }, 4000);
            }
        } catch (err) {
            console.error("Delete error:", err);
            toastMessage = "Terjadi kesalahan saat menghapus data";
            toastType = "error";
            showToast = true;
            setTimeout(() => {
                showToast = false;
            }, 4000);
        } finally {
            isDeleting = false;
        }
    }

    async function executeExport(selectedType: string) {
        if (!selectedType) return;

        try {
            isExporting = true;
            toastMessage = "Menyiapkan file Excel...";
            toastType = "info";
            showToast = true;

            const queryParams = `?jenis_situs=${encodeURIComponent(selectedType)}`;
            const response = await apiFetch(`/api/situs/export${queryParams}`);

            if (!response.ok) {
                const errText = await response.text();
                throw new Error(errText || "Gagal mengekspor data");
            }

            const blob = await response.blob();
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement("a");
            a.href = url;

            a.download = `Data_Situs_${selectedType}.xlsx`;
            document.body.appendChild(a);
            a.click();

            a.remove();
            window.URL.revokeObjectURL(url);

            showExportModal = false;
            showToast = false;
        } catch (err: any) {
            console.error("Export error:", err);
            toastMessage = "Gagal mengunduh Excel: " + err.message;
            toastType = "error";
            showToast = true;
            setTimeout(() => (showToast = false), 4000);
        } finally {
            isExporting = false;
        }
    }
</script>

<svelte:head>
    <title>Data Situs | SiCemas</title>
    <meta name="description" content="Sistem Informasi Keagamaan KUA Ciemas" />
</svelte:head>

<div class="w-full">
    <Toast
        show={showToast}
        message={toastMessage}
        type={toastType}
        onclose={() => (showToast = false)}
    />

    <ConfirmModal
        show={showDeleteModal}
        title="Hapus Data Situs"
        message="Apakah Anda yakin ingin menghapus {siteToDelete?.nama}? Data yang dihapus tidak dapat dikembalikan."
        confirmText="Hapus"
        cancelText="Batal"
        isProcessing={isDeleting}
        onConfirm={handleDelete}
        onCancel={() => (showDeleteModal = false)}
        type="danger"
    />

    <ExportSitusModal
        show={showExportModal}
        {jenisSitus}
        {isExporting}
        onClose={() => (showExportModal = false)}
        onExport={executeExport}
    />

    <div
        class="mb-8 flex flex-col gap-6 md:flex-row md:items-end md:justify-between"
    >
        <div>
            <div class="flex items-center gap-3 mb-2">
                <div
                    class="w-1.5 h-8 rounded-full bg-linear-to-b from-primary to-secondary"
                ></div>
                <h1
                    class="text-3xl font-extrabold tracking-tight text-base-content"
                >
                    Daftar Situs Keagamaan
                </h1>
            </div>
            <p class="ml-5 font-medium text-base-content/60 text-sm">
                Kelola dan pantau situs-situs keagamaan di wilayah Kecamatan
                Ciemas
            </p>
        </div>

        <div class="flex items-center gap-3">
            <button
                type="button"
                title="Refresh Data"
                onclick={handleRefresh}
                disabled={isRefreshing}
                class="btn btn-circle btn-ghost bg-base-200/50 hover:bg-base-300 transition-colors"
            >
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-5 w-5 text-base-content/70 {isRefreshing
                        ? 'animate-spin text-primary'
                        : ''}"
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
            </button>

            {#if auth.hasAllPermissions(["situs:create"])}
                <button
                    type="button"
                    class="btn btn-primary rounded-xl shadow-lg shadow-primary/20 transition-all hover:-translate-y-0.5 hover:shadow-primary/30"
                    onclick={() => goto("/dashboard/situs/tambah")}
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
                            d="M12 6v6m0 0v6m0-6h6m-6 0H6"
                        />
                    </svg>
                    Tambah Situs
                </button>
            {/if}
        </div>
    </div>

    <div
        class="mb-6 rounded-2xl bg-base-100 p-4 sm:p-5 shadow-sm border border-base-200"
    >
        <div
            class="flex flex-col lg:flex-row gap-4 lg:items-center justify-between"
        >
            <div
                class="flex flex-col sm:flex-row gap-3 w-full lg:w-auto flex-1"
            >
                <div class="relative w-full sm:max-w-xs md:max-w-md">
                    <div
                        class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-4"
                    >
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="h-5 w-5 text-base-content/40"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                            />
                        </svg>
                    </div>
                    <input
                        type="text"
                        placeholder="Cari situs, ID, lokasi..."
                        class="input input-bordered w-full rounded-xl bg-base-200/50 pl-11 focus:bg-base-100 transition-colors"
                        bind:value={searchTerm}
                    />
                </div>

                <select
                    class="select select-bordered rounded-xl bg-base-200/50 focus:bg-base-100 transition-colors w-full sm:w-auto"
                    bind:value={statusFilter}
                >
                    <option value="all">Semua Status</option>
                    <option value="menunggu">Menunggu</option>
                    <option value="terverifikasi">Terverifikasi</option>
                    <option value="ditolak">Ditolak</option>
                </select>

                <select
                    class="select select-bordered rounded-xl bg-base-200/50 focus:bg-base-100 transition-colors w-full sm:w-auto"
                    bind:value={typeFilter}
                >
                    <option value="all">Semua Jenis</option>
                    {#each jenisSitus as jenis (jenis.id)}
                        <option value={jenis.nama}>{jenis.nama}</option>
                    {/each}
                </select>
            </div>

            <div
                class="flex items-center justify-between lg:justify-end gap-4 w-full lg:w-auto border-t border-base-200 pt-4 lg:border-t-0 lg:pt-0"
            >
                <div
                    class="text-sm font-medium text-base-content/60 whitespace-nowrap"
                >
                    Total <span class="text-base-content font-bold"
                        >{filteredSites.length}</span
                    > Data
                </div>
                {#if auth.hasAllPermissions(["situs:export"])}
                    <button
                        type="button"
                        onclick={() => (showExportModal = true)}
                        class="btn btn-outline btn-success rounded-xl hover:text-white transition-colors"
                        title="Export data to Excel"
                    >
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="h-4 w-4"
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
                        <span class="hidden sm:inline">Export</span> Excel
                    </button>
                {/if}
            </div>
        </div>
    </div>

    <div
        class="overflow-hidden rounded-2xl border border-base-200 bg-base-100 shadow-xl shadow-base-200/40"
    >
        {#if isRefreshing}
            <div
                class="flex flex-col items-center justify-center py-20 px-4 text-center"
            >
                <span class="loading loading-lg loading-spinner text-primary"
                ></span>
                <p class="mt-4 font-medium text-base-content/70">
                    Menyegarkan data dari server...
                </p>
            </div>
        {:else if filteredSites.length === 0}
            <div
                class="flex flex-col items-center justify-center py-20 px-4 text-center"
            >
                <div
                    class="mb-4 flex h-20 w-20 items-center justify-center rounded-full bg-base-200/50"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="h-10 w-10 text-base-content/30"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                        />
                    </svg>
                </div>
                <h3 class="text-lg font-bold text-base-content">
                    Pencarian Tidak Ditemukan
                </h3>
                <p class="mt-1 text-sm text-base-content/60">
                    Coba sesuaikan kata kunci atau filter status/jenis yang
                    lain.
                </p>
            </div>
        {:else}
            <div class="overflow-x-auto">
                <table class="table w-full">
                    <thead class="bg-base-200/50 text-base-content/70">
                        <tr
                            class="text-xs uppercase tracking-wider font-bold border-b-2 border-base-200"
                        >
                            <th class="py-4">Info Situs</th>
                            <th>Kategori & Lokasi</th>
                            <th>Pendata</th>
                            <th class="text-center">Status</th>
                            <th class="text-center">Aksi</th>
                        </tr>
                    </thead>
                    <tbody class="divide-y divide-base-200">
                        {#each paginatedSites as site (site.id)}
                            <tr class="transition-colors hover:bg-base-200/30">
                                <td class="py-3">
                                    <div
                                        class="font-extrabold text-base-content text-sm mb-0.5"
                                    >
                                        {site.nama}
                                    </div>
                                    <div
                                        class="flex items-center gap-1.5 text-xs font-mono text-base-content/50"
                                    >
                                        <svg
                                            xmlns="http://www.w3.org/2000/svg"
                                            class="h-3 w-3"
                                            fill="none"
                                            viewBox="0 0 24 24"
                                            stroke="currentColor"
                                        >
                                            <path
                                                stroke-linecap="round"
                                                stroke-linejoin="round"
                                                stroke-width="2"
                                                d="M7 20l4-16m2 16l4-16M6 9h14M4 15h14"
                                            />
                                        </svg>
                                        {site.situs_id || "Belum ada ID"}
                                    </div>
                                </td>
                                <td>
                                    <div class="mb-1">
                                        <span
                                            class="badge border-primary/20 bg-primary/10 text-primary text-xs font-semibold"
                                        >
                                            {site.jenis}
                                        </span>
                                    </div>
                                    <div
                                        class="text-xs text-base-content/70 flex items-center gap-1"
                                    >
                                        <svg
                                            xmlns="http://www.w3.org/2000/svg"
                                            class="h-3 w-3"
                                            fill="none"
                                            viewBox="0 0 24 24"
                                            stroke="currentColor"
                                        >
                                            <path
                                                stroke-linecap="round"
                                                stroke-linejoin="round"
                                                stroke-width="2"
                                                d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"
                                            />
                                            <path
                                                stroke-linecap="round"
                                                stroke-linejoin="round"
                                                stroke-width="2"
                                                d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"
                                            />
                                        </svg>
                                        <span
                                            class="truncate max-w-37.5"
                                            title={site.lokasi}
                                            >{site.lokasi}</span
                                        >
                                    </div>
                                </td>
                                <td>
                                    <div
                                        class="text-sm font-medium text-base-content/80"
                                    >
                                        {site.pendata}
                                    </div>
                                    <div
                                        class="text-[10px] text-base-content/40 mt-0.5"
                                    >
                                        Update: {formatDate(site.updated_at)}
                                    </div>
                                </td>
                                <td class="text-center">
                                    <div
                                        class={`badge badge-sm font-medium px-2.5 py-2.5 ${
                                            site.status === "terverifikasi"
                                                ? "bg-success/10 text-success border-success/20"
                                                : site.status === "menunggu"
                                                  ? "bg-warning/10 text-warning border-warning/20"
                                                  : "bg-error/10 text-error border-error/20"
                                        }`}
                                    >
                                        <div
                                            class={`w-1.5 h-1.5 rounded-full mr-1.5 ${
                                                site.status === "terverifikasi"
                                                    ? "bg-success"
                                                    : site.status === "menunggu"
                                                      ? "bg-warning"
                                                      : "bg-error"
                                            }`}
                                        ></div>
                                        {site.status.charAt(0).toUpperCase() +
                                            site.status.slice(1)}
                                    </div>
                                </td>
                                <td>
                                    <div class="flex justify-center">
                                        <SitusTableActions
                                            id={site.id}
                                            status={site.status}
                                            onDelete={() => confirmDelete(site)}
                                        />
                                    </div>
                                </td>
                            </tr>
                        {/each}
                    </tbody>
                </table>
            </div>

            <div class="border-t border-base-200 bg-base-100/50 p-4">
                <Pagination
                    bind:currentPage
                    bind:itemsPerPage
                    totalItems={filteredSites.length}
                />
            </div>
        {/if}
    </div>
</div>
