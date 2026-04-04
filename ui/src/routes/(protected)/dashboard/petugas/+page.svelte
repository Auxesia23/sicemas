<script lang="ts">
    import type { PageData } from "./$types";
    import type { UserResponse } from "$lib/schemas/user";
    import { UserInputSchema } from "$lib/schemas/user";

    import * as v from "valibot";
    import { auth } from "$lib/states/auth.svelte";
    import { apiFetch } from "$lib/api";
    import { invalidate } from "$app/navigation";

    import UserTableActions from "$lib/components/features/users/UserTableActions.svelte";
    import UserFormModal from "$lib/components/features/users/UserFormModal.svelte";
    import ConfirmModal from "$lib/components/ui/ConfirmModal.svelte";
    import Toast from "$lib/components/ui/Toast.svelte";
    import Pagination from "$lib/components/ui/Pagination.svelte";

    let { data }: { data: PageData } = $props();
    let { users: staffs, roles } = $derived(data);

    let isModalOpen = $state(false);
    let isEditMode = $state(false);
    let editingUserId = $state<string | null>(null);
    let formLoading = $state(false);

    let modalError = $state("");

    let errors = $state<{
        nip?: string;
        nama_lengkap?: string;
        email?: string;
        nomor_telepon?: string;
        jabatan?: string;
        unit_kerja?: string;
        peran?: string;
    }>({});

    let currentPage = $state(1);
    let itemsPerPage = $state(10);
    let searchTerm = $state("");

    let showDeleteModal = $state(false);
    let userToDelete = $state<UserResponse | null>(null);
    let isDeleting = $state(false);
    let showToast = $state(false);
    let toastMessage = $state("");
    let toastType = $state<"success" | "error" | "warning" | "info">("success");

    let formData = $state<Record<string, any>>({
        nip: "",
        nama_lengkap: "",
        email: "",
        nomor_telepon: "",
        jabatan: "",
        unit_kerja: "",
        peran: "",
    });

    let filteredStaffs = $derived(
        staffs.filter((staff: UserResponse) => {
            const searchLower = searchTerm.toLowerCase();
            return (
                staff.nama_lengkap.toLowerCase().includes(searchLower) ||
                staff.nip.toLowerCase().includes(searchLower) ||
                staff.jabatan.toLowerCase().includes(searchLower) ||
                staff.peran.toLowerCase().includes(searchLower)
            );
        }),
    );

    let paginatedStaffs = $derived(
        filteredStaffs.slice(
            (currentPage - 1) * itemsPerPage,
            currentPage * itemsPerPage,
        ),
    );

    $effect(() => {
        searchTerm;
        itemsPerPage;
        currentPage = 1;
    });

    const openModal = () => {
        isModalOpen = true;
        isEditMode = false;
        editingUserId = null;
        modalError = "";
        errors = {};
        formData = {
            nip: "",
            nama_lengkap: "",
            email: "",
            nomor_telepon: "",
            jabatan: "",
            unit_kerja: "",
            peran: "",
        };
    };

    const openEditModal = (user: UserResponse) => {
        isModalOpen = true;
        isEditMode = true;
        editingUserId = user.id;
        modalError = "";
        errors = {};
        formData = {
            nip: user.nip || "",
            nama_lengkap: user.nama_lengkap || "",
            email: user.email || "",
            nomor_telepon: user.nomor_telepon || "",
            jabatan: user.jabatan || "",
            unit_kerja: user.unit_kerja || "",
            peran: user.peran || "",
        };
    };

    const closeModal = () => {
        isModalOpen = false;
        errors = {};
    };

    const handleSubmit = async () => {
        modalError = "";
        errors = {};

        const parsed = v.safeParse(UserInputSchema, formData);

        if (!parsed.success) {
            formLoading = false;
            const flat = v.flatten(parsed.issues);
            errors = {
                nip: flat.nested?.nip?.[0],
                nama_lengkap: flat.nested?.nama_lengkap?.[0],
                email: flat.nested?.email?.[0],
                nomor_telepon: flat.nested?.nomor_telepon?.[0],
                jabatan: flat.nested?.jabatan?.[0],
                unit_kerja: flat.nested?.unit_kerja?.[0],
                peran: flat.nested?.peran?.[0],
            };
            return;
        }

        formLoading = true;
        const payload = parsed.output;

        try {
            let response;
            if (isEditMode) {
                response = await apiFetch(`/api/users/${editingUserId}`, {
                    method: "PUT",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify(payload),
                });
            } else {
                response = await apiFetch("/api/users", {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify(payload),
                });
            }

            if (!response.ok) {
                const textError = await response.text();
                throw new Error(textError || "Gagal menyimpan data petugas");
            }

            closeModal();
            toastMessage = isEditMode
                ? "Data petugas berhasil diperbarui"
                : "Petugas berhasil ditambahkan";
            toastType = "success";
            showToast = true;

            await invalidate("data:users");

            setTimeout(() => {
                showToast = false;
            }, 3000);
        } catch (error: any) {
            modalError = error.message;
        } finally {
            formLoading = false;
        }
    };

    const confirmDelete = (user: UserResponse) => {
        userToDelete = user;
        showDeleteModal = true;
    };

    const deleteUser = async () => {
        if (!userToDelete) return;
        isDeleting = true;

        try {
            const response = await apiFetch(`/api/users/${userToDelete.id}`, {
                method: "DELETE",
            });

            if (!response.ok) {
                const textError = await response.text();
                throw new Error(textError || "Gagal menghapus petugas");
            }

            showDeleteModal = false;
            toastMessage = "Petugas berhasil dihapus";
            toastType = "success";
            showToast = true;

            await invalidate("data:users");

            setTimeout(() => {
                showToast = false;
            }, 3000);
        } catch (error: any) {
            toastMessage = error.message;
            toastType = "error";
            showToast = true;
            setTimeout(() => {
                showToast = false;
            }, 4000);
        } finally {
            isDeleting = false;
            userToDelete = null;
        }
    };
</script>

<svelte:head>
    <title>Petugas | SiCemas</title>
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
        title="Hapus Petugas"
        message="Apakah Anda yakin ingin menghapus {userToDelete?.nama_lengkap}? Tindakan ini tidak dapat dibatalkan."
        type="danger"
        isProcessing={isDeleting}
        onConfirm={deleteUser}
        onCancel={() => (showDeleteModal = false)}
    />

    <UserFormModal
        show={isModalOpen}
        {isEditMode}
        isLoading={formLoading}
        bind:formData
        {roles}
        {modalError}
        bind:errors
        onSubmit={handleSubmit}
        onCancel={closeModal}
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
                    Data Petugas
                </h1>
            </div>
            <p class="ml-5 font-medium text-base-content/60 text-sm">
                Manajemen informasi personil KUA Kecamatan Ciemas
            </p>
        </div>

        {#if auth.hasAllPermissions(["user:create"])}
            <button
                type="button"
                class="btn btn-primary rounded-xl shadow-lg shadow-primary/20 transition-all hover:-translate-y-0.5 hover:shadow-primary/30"
                onclick={openModal}
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
                Tambah Petugas
            </button>
        {/if}
    </div>

    <div
        class="mb-6 rounded-2xl bg-base-100 p-4 shadow-sm border border-base-200 flex flex-col gap-4 sm:flex-row sm:items-center"
    >
        <div class="relative w-full sm:max-w-md">
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
                placeholder="Cari nama, NIP, atau jabatan..."
                class="input input-bordered w-full rounded-xl bg-base-200/50 pl-11 focus:bg-base-100 transition-colors"
                bind:value={searchTerm}
            />
        </div>

        <div class="flex-1"></div>
        <div class="text-sm font-medium text-base-content/50">
            Total {filteredStaffs.length} Petugas
        </div>
    </div>

    <div
        class="overflow-hidden rounded-2xl border border-base-200 bg-base-100 shadow-xl shadow-base-200/40"
    >
        {#if filteredStaffs.length === 0}
            <div
                class="flex flex-col items-center justify-center py-16 px-4 text-center"
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
                            d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
                        />
                    </svg>
                </div>
                <h3 class="text-lg font-bold text-base-content">
                    Tidak ada petugas ditemukan
                </h3>
                <p class="mt-1 text-sm text-base-content/60">
                    Coba ubah kata kunci pencarian atau tambah petugas baru.
                </p>
            </div>
        {:else}
            <div class="overflow-x-auto">
                <table class="table w-full">
                    <thead class="bg-base-200/50 text-base-content/70">
                        <tr
                            class="text-xs uppercase tracking-wider font-bold border-b-2 border-base-200"
                        >
                            <th class="py-4">Profil Pegawai</th>
                            <th>Jabatan & Unit</th>
                            <th>Kontak</th>
                            <th class="text-center">Peran</th>
                            <th class="text-center w-24">Aksi</th>
                        </tr>
                    </thead>
                    <tbody class="divide-y divide-base-200">
                        {#each paginatedStaffs as staff (staff.id)}
                            <tr class="transition-colors hover:bg-base-200/30">
                                <td class="py-3">
                                    <div class="flex items-center gap-3">
                                        <div class="avatar placeholder">
                                            <div
                                                class="w-10 h-10 flex items-center justify-center rounded-full bg-linear-to-br from-primary/20 to-secondary/20 text-primary font-bold"
                                            >
                                                <span
                                                    >{staff.nama_lengkap
                                                        .charAt(0)
                                                        .toUpperCase()}</span
                                                >
                                            </div>
                                        </div>
                                        <div>
                                            <div
                                                class="font-bold text-base-content"
                                            >
                                                {staff.nama_lengkap}
                                            </div>
                                            <div
                                                class="text-xs font-mono text-base-content/60 mt-0.5"
                                            >
                                                NIP: {staff.nip}
                                            </div>
                                        </div>
                                    </div>
                                </td>
                                <td>
                                    <div
                                        class="font-medium text-base-content/80"
                                    >
                                        {staff.jabatan}
                                    </div>
                                    <div
                                        class="text-xs text-base-content/50 mt-0.5"
                                    >
                                        {staff.unit_kerja || "-"}
                                    </div>
                                </td>
                                <td>
                                    <div class="text-sm text-base-content/80">
                                        {staff.nomor_telepon}
                                    </div>
                                    <div
                                        class="text-xs text-base-content/50 mt-0.5"
                                    >
                                        {staff.email || "-"}
                                    </div>
                                </td>
                                <td class="text-center">
                                    <span
                                        class="badge border-primary/20 bg-primary/10 px-3 py-3 font-semibold text-primary"
                                    >
                                        {staff.peran.toUpperCase()}
                                    </span>
                                </td>
                                <td>
                                    <div class="flex justify-center">
                                        <UserTableActions
                                            onUpdate={() =>
                                                openEditModal(staff)}
                                            onDelete={() =>
                                                confirmDelete(staff)}
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
                    totalItems={filteredStaffs.length}
                />
            </div>
        {/if}
    </div>
</div>
