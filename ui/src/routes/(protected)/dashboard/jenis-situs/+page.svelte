<script lang="ts">
    import { auth } from "$lib/states/auth.svelte";
    import { apiFetch } from "$lib/api";
    import { invalidate } from "$app/navigation";
    import ConfirmModal from "$lib/components/ui/ConfirmModal.svelte";
    import Toast from "$lib/components/ui/Toast.svelte";

    import type { JenisSitus } from "$lib/schemas/sites";
    import type { PageData } from "./$types";

    import * as v from "valibot";
    import { JenisSitusInputSchema } from "$lib/schemas/sites";

    let { data }: { data: PageData } = $props();

    const colorMap = [
        "primary",
        "secondary",
        "accent",
        "info",
        "success",
        "warning",
        "error",
    ];

    let siteTypes = $derived(
        data.jenisSitus.map((item, index) => ({
            ...item,
            color: colorMap[index % colorMap.length],
        })),
    );

    let modalOpen = $state(false);
    let isEditMode = $state(false);
    let editId = $state<string | null>(null);
    let formNama = $state("");
    let formDeskripsi = $state("");
    let formLoading = $state(false);

    let errors = $state<{ nama?: string; deskripsi?: string }>({});

    let showDeleteModal = $state(false);
    let itemToDelete = $state<{ id: string; name: string } | null>(null);
    let isDeleting = $state(false);

    let showToast = $state(false);
    let toastMessage = $state("");
    let toastType = $state<"success" | "error" | "warning" | "info">("success");

    function openModal(item: JenisSitus | null = null) {
        modalOpen = true;
        formLoading = false;
        errors = {};

        if (item) {
            isEditMode = true;
            editId = item.id;
            formNama = item.nama;
            formDeskripsi = item.deskripsi;
        } else {
            isEditMode = false;
            editId = null;
            formNama = "";
            formDeskripsi = "";
        }
    }

    function closeModal() {
        modalOpen = false;
        errors = {};
    }

    function clearError(field: "nama" | "deskripsi") {
        if (errors[field]) {
            errors = { ...errors, [field]: undefined };
        }
    }

    async function submitForm() {
        formLoading = true;
        errors = {};

        const rawPayload = {
            nama: formNama,
            deskripsi: formDeskripsi,
        };

        const parsed = v.safeParse(JenisSitusInputSchema, rawPayload);

        if (!parsed.success) {
            formLoading = false;
            const flat = v.flatten(parsed.issues);
            errors = {
                nama: flat.nested?.nama?.[0],
                deskripsi: flat.nested?.deskripsi?.[0],
            };
            return;
        }

        try {
            const payload = parsed.output;
            let response;

            if (isEditMode) {
                response = await apiFetch(`/api/jenis-situs/${editId}`, {
                    method: "PUT",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify(payload),
                });
            } else {
                response = await apiFetch("/api/jenis-situs", {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify(payload),
                });
            }

            if (!response.ok) {
                const text = await response.text();
                throw new Error(
                    text ||
                        `Gagal ${isEditMode ? "mengubah" : "menambah"} jenis situs`,
                );
            }

            toastMessage = `Jenis situs berhasil ${isEditMode ? "diperbarui" : "ditambahkan"}`;
            toastType = "success";
            showToast = true;

            await invalidate("data:jenis-situs");

            closeModal();
            setTimeout(() => {
                showToast = false;
            }, 3000);
        } catch (e: any) {
            toastMessage =
                e.message ||
                `Gagal ${isEditMode ? "mengubah" : "menambah"} jenis situs`;
            toastType = "error";
            showToast = true;
            setTimeout(() => {
                showToast = false;
            }, 4000);
        } finally {
            formLoading = false;
        }
    }

    function confirmDeleteItem(id: string, name: string) {
        itemToDelete = { id, name };
        showDeleteModal = true;
    }

    async function handleDelete() {
        if (!itemToDelete) return;

        isDeleting = true;

        try {
            const response = await apiFetch(
                `/api/jenis-situs/${itemToDelete.id}`,
                {
                    method: "DELETE",
                },
            );

            if (!response.ok) {
                const text = await response.text();
                throw new Error(text || "Gagal menghapus jenis situs");
            }

            toastMessage = "Jenis situs berhasil dihapus";
            toastType = "success";
            showToast = true;

            await invalidate("data:jenis-situs");

            showDeleteModal = false;
            setTimeout(() => {
                showToast = false;
            }, 3000);
        } catch (e: any) {
            toastMessage = e.message || "Terjadi kesalahan saat menghapus data";
            toastType = "error";
            showToast = true;
            setTimeout(() => {
                showToast = false;
            }, 4000);
        } finally {
            showDeleteModal = false;
            isDeleting = false;
            itemToDelete = null;
        }
    }
</script>

<svelte:head>
    <title>Jenis Situs | SiCemas</title>
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
        title="Hapus Jenis Situs"
        message="Apakah Anda yakin ingin menghapus jenis situs '{itemToDelete?.name}'? Tindakan ini tidak dapat dibatalkan."
        type="danger"
        confirmText="Hapus"
        cancelText="Batal"
        isProcessing={isDeleting}
        onConfirm={handleDelete}
        onCancel={() => (showDeleteModal = false)}
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
                    Jenis Situs Keagamaan
                </h1>
            </div>
            <p class="ml-5 font-medium text-base-content/60 text-sm">
                Kelola kategori jenis-jenis situs keagamaan yang terdaftar
            </p>
        </div>

        {#if auth.hasAllPermissions(["jenis-situs:create"])}
            <button
                type="button"
                class="btn btn-primary rounded-xl shadow-lg shadow-primary/20 transition-all hover:-translate-y-0.5 hover:shadow-primary/30"
                onclick={() => openModal()}
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
                Tambah Jenis Situs
            </button>
        {/if}
    </div>

    {#if siteTypes.length === 0}
        <div
            class="flex flex-col items-center justify-center py-20 px-4 text-center rounded-2xl border border-base-200 bg-base-100 shadow-xl shadow-base-200/40"
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
                        d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 002-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"
                    />
                </svg>
            </div>
            <h3 class="text-lg font-bold text-base-content">
                Belum ada jenis situs
            </h3>
            <p class="mt-1 text-sm text-base-content/60">
                Silakan tambah kategori jenis situs keagamaan pertama Anda.
            </p>
        </div>
    {:else}
        <div class="grid grid-cols-1 gap-6 md:grid-cols-2 xl:grid-cols-3">
            {#each siteTypes as type (type.id)}
                <div
                    class="card rounded-2xl border border-base-200 bg-base-100 shadow-xl shadow-base-200/30 transition-all duration-300 hover:-translate-y-1 hover:shadow-2xl hover:shadow-base-200/50"
                >
                    <div class="card-body p-6">
                        <div class="flex items-start justify-between gap-4">
                            <div>
                                <h3
                                    class="text-lg font-extrabold text-base-content mb-1"
                                >
                                    {type.nama}
                                </h3>
                                <p
                                    class="text-sm font-medium text-base-content/60 line-clamp-2"
                                >
                                    {type.deskripsi || "Tidak ada deskripsi"}
                                </p>
                            </div>
                            <div
                                class="badge border-{type.color}/20 bg-{type.color}/10 text-{type.color} font-bold px-3 py-3 rounded-lg whitespace-nowrap"
                            >
                                {type.jumlah_situs} situs
                            </div>
                        </div>

                        <div
                            class="mt-6 flex items-center justify-end gap-2 border-t border-base-200 pt-4"
                        >
                            {#if auth.hasAllPermissions(["jenis-situs:update"])}
                                <button
                                    type="button"
                                    class="btn btn-ghost btn-sm rounded-lg text-base-content/70 hover:bg-primary/10 hover:text-primary transition-colors"
                                    onclick={() => openModal(type)}
                                    title="Edit Jenis Situs"
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
                                            d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                                        />
                                    </svg>
                                    Edit
                                </button>
                            {/if}
                            {#if auth.hasAllPermissions(["jenis-situs:delete"])}
                                <button
                                    type="button"
                                    class="btn btn-ghost btn-sm rounded-lg text-base-content/70 hover:bg-error/10 hover:text-error transition-colors"
                                    onclick={() =>
                                        confirmDeleteItem(type.id, type.nama)}
                                    title="Hapus Jenis Situs"
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
                                            d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                                        />
                                    </svg>
                                    Hapus
                                </button>
                            {/if}
                        </div>
                    </div>
                </div>
            {/each}
        </div>
    {/if}
</div>

<dialog class="modal modal-bottom sm:modal-middle" open={modalOpen}>
    <div
        class="modal-box p-0 max-w-lg rounded-2xl shadow-2xl overflow-hidden bg-base-100"
    >
        <div
            class="bg-base-200/40 px-6 py-5 border-b border-base-200 sm:px-8 flex justify-between items-center"
        >
            <div class="flex items-center gap-4">
                <div
                    class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-linear-to-br from-primary/20 to-secondary/20 text-primary shadow-sm"
                >
                    {#if isEditMode}
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
                                d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                            />
                        </svg>
                    {:else}
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
                    {/if}
                </div>
                <h3 class="text-xl font-extrabold text-base-content">
                    {isEditMode ? "Edit Jenis Situs" : "Tambah Jenis Situs"}
                </h3>
            </div>
            <button
                type="button"
                class="btn btn-circle btn-ghost btn-sm"
                onclick={closeModal}
                disabled={formLoading}
            >
                ✕
            </button>
        </div>

        <div class="px-6 py-6 sm:px-8 space-y-5">
            <div class="form-control">
                <label class="label pb-1.5 pt-0" for="formNama">
                    <span class="label-text font-bold text-base-content/80">
                        Nama Jenis Situs <span class="text-error">*</span>
                    </span>
                </label>
                <input
                    id="formNama"
                    type="text"
                    class="input input-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors {errors.nama
                        ? 'input-error focus:border-error'
                        : 'focus:border-primary'}"
                    bind:value={formNama}
                    placeholder="Contoh: Masjid, Pondok Pesantren..."
                    oninput={() => clearError("nama")}
                    disabled={formLoading}
                />
                {#if errors.nama}
                    <span class="mt-1.5 text-xs font-medium text-error"
                        >{errors.nama}</span
                    >
                {/if}
            </div>

            <div class="form-control">
                <label class="label pb-1.5 pt-0" for="formDeskripsi">
                    <span class="label-text font-bold text-base-content/80">
                        Deskripsi<span class="text-error">*</span>
                    </span>
                </label>
                <textarea
                    id="formDeskripsi"
                    class="textarea textarea-bordered h-28 w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors {errors.deskripsi
                        ? 'textarea-error focus:border-error'
                        : 'focus:border-primary'}"
                    bind:value={formDeskripsi}
                    placeholder="Tuliskan keterangan singkat mengenai jenis situs ini..."
                    oninput={() => clearError("deskripsi")}
                    disabled={formLoading}
                ></textarea>
                {#if errors.deskripsi}
                    <span class="mt-1.5 text-xs font-medium text-error"
                        >{errors.deskripsi}</span
                    >
                {/if}
            </div>
        </div>

        <div
            class="bg-base-200/40 px-6 py-5 border-t border-base-200 sm:px-8 flex flex-col-reverse sm:flex-row justify-end gap-3"
        >
            <button
                type="button"
                class="btn btn-ghost rounded-xl hover:bg-base-300 w-full sm:w-auto"
                onclick={closeModal}
                disabled={formLoading}
            >
                Batal
            </button>
            <button
                type="button"
                class="btn btn-primary rounded-xl px-8 shadow-lg shadow-primary/20 transition-transform hover:-translate-y-0.5 w-full sm:w-auto text-primary-content"
                onclick={submitForm}
                disabled={formLoading}
            >
                {#if formLoading}
                    <span class="loading loading-sm loading-spinner"></span>
                    Menyimpan...
                {:else}
                    {isEditMode ? "Perbarui" : "Simpan"}
                {/if}
            </button>
        </div>
    </div>

    <button
        type="button"
        class="modal-backdrop bg-base-content/20 backdrop-blur-sm"
        onclick={closeModal}
        aria-label="Tutup"
        disabled={formLoading}
    ></button>
</dialog>
