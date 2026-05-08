<script lang="ts">
    import { apiFetch } from "$lib/api";
    import { auth } from "$lib/states/auth.svelte";
    import { invalidate } from "$app/navigation";

    import Toast from "$lib/components/ui/Toast.svelte";
    import ConfirmModal from "$lib/components/ui/ConfirmModal.svelte";
    import RoleBadge from "$lib/components/features/users/RoleBadge.svelte";
    import AddRoleForm from "$lib/components/features/users/AddRoleForm.svelte";

    import type { PageData } from "./$types";
    import type { RoleResponse } from "$lib/schemas/user";

    let { data }: { data: PageData } = $props();

    let { roles, policies } = $derived(data);
    let isAdding = $state(false);

    const resources = [
        { name: "dashboard", actions: ["read"] },
        { name: "user", actions: ["create", "read", "update", "delete"] },
        {
            name: "situs",
            actions: [
                "create",
                "read_all",
                "read_own",
                "update",
                "delete",
                "verify",
                "export",
            ],
        },
        {
            name: "jenis-situs",
            actions: ["create", "read", "update", "delete"],
        },
        { name: "role", actions: ["create", "read", "delete"] },
        { name: "policy", actions: ["create", "read", "update", "delete"] },
        { name: "activity", actions: ["read"] },
    ];

    function getActionLabel(action: string): string {
        const labels: Record<string, string> = {
            create: "Create",
            read: "Read",
            read_all: "Read All",
            read_own: "Read Own",
            update: "Update",
            update_all: "Update All",
            update_own: "Update Own",
            delete: "Delete",
            verify: "Verify",
            export: "Export",
        };
        return labels[action] || action;
    }

    let updatingCells = $state<Record<string, boolean>>({});

    let showDeleteModal = $state(false);
    let roleToDelete = $state<RoleResponse | null>(null);
    let isDeleting = $state(false);
    let showToast = $state(false);
    let toastMessage = $state("");
    let toastType = $state<"success" | "error" | "warning" | "info">("success");

    function getPolicyKey(role: string, resource: string, action: string) {
        return `${role}-${resource}-${action}`;
    }

    function hasPolicy(role: string, resource: string, action: string) {
        return policies.some(
            (p) =>
                p.subject === role &&
                p.object === resource &&
                p.action === action,
        );
    }

    function isCellUpdating(role: string, resource: string, action: string) {
        return !!updatingCells[getPolicyKey(role, resource, action)];
    }

    async function handleAddNewRoleWrapper(name: string) {
        if (!name) return;
        isAdding = true;

        try {
            const res = await apiFetch("/api/roles", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ name }),
            });

            if (res.ok) {
                toastMessage = "Role berhasil ditambahkan";
                toastType = "success";
                showToast = true;

                await invalidate("data:permissions");

                setTimeout(() => {
                    showToast = false;
                }, 3000);
            } else {
                const errorText = await res.text();
                throw new Error(errorText || "Gagal menambahkan role");
            }
        } catch (err: any) {
            toastMessage =
                err.message || "Terjadi kesalahan saat menambahkan role";
            toastType = "error";
            showToast = true;
            setTimeout(() => {
                showToast = false;
            }, 4000);
        } finally {
            isAdding = false;
        }
    }

    function confirmDeleteRole(role: RoleResponse) {
        roleToDelete = role;
        showDeleteModal = true;
    }

    async function deleteRole() {
        if (!roleToDelete) return;
        isDeleting = true;

        try {
            const res = await apiFetch(`/api/roles/${roleToDelete.id}`, {
                method: "DELETE",
            });

            if (res.ok) {
                showDeleteModal = false;
                toastMessage = "Role berhasil dihapus";
                toastType = "success";
                showToast = true;

                await invalidate("data:permissions");

                setTimeout(() => {
                    showToast = false;
                }, 3000);
            } else {
                const errorText = await res.text();
                throw new Error(
                    errorText || `Gagal menghapus role ${roleToDelete.name}`,
                );
            }
        } catch (err: any) {
            toastMessage =
                err.message || "Terjadi kesalahan saat menghapus role";
            toastType = "error";
            showToast = true;
            setTimeout(() => {
                showToast = false;
            }, 4000);
        } finally {
            isDeleting = false;
            roleToDelete = null;
        }
    }

    async function handleCheckboxChange(
        role: string,
        resource: string,
        action: string,
        checked: boolean,
    ) {
        const key = getPolicyKey(role, resource, action);
        updatingCells[key] = true;

        try {
            let res;
            const payload = {
                subject: role,
                object: resource,
                action: action,
            };

            if (checked) {
                res = await apiFetch("/api/policies", {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify(payload),
                });
            } else {
                res = await apiFetch("/api/policies", {
                    method: "DELETE",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify(payload),
                });
            }

            if (res.ok) {
                await invalidate("data:permissions");
                await auth.fetchUser();
            } else {
                const errorText = await res.text();
                throw new Error(
                    errorText || `Gagal mengubah hak akses ${role}`,
                );
            }
        } catch (err: any) {
            toastMessage =
                err.message || "Terjadi kesalahan saat mengubah hak akses";
            toastType = "error";
            showToast = true;
            setTimeout(() => {
                showToast = false;
            }, 4000);
        } finally {
            updatingCells[key] = false;
        }
    }
</script>

<svelte:head>
    <title>Policy & Role | SiCemas</title>
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
        title="Hapus Role"
        message="Apakah Anda yakin ingin menghapus role '{roleToDelete?.name}'? Semua hak akses terkait juga akan ikut terhapus permanen."
        type="danger"
        confirmText="Hapus"
        cancelText="Batal"
        isProcessing={isDeleting}
        onConfirm={deleteRole}
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
                    class="text-3xl font-extrabold tracking-tight text-base-content leading-none"
                >
                    Manajemen Hak Akses
                </h1>
            </div>
            <p class="ml-5 font-medium text-base-content/60 text-sm">
                Kelola hak akses (policy matrix) untuk setiap peran dan sumber
                daya sistem
            </p>
        </div>
    </div>

    <div class="mb-8 grid grid-cols-1 gap-4 sm:grid-cols-3">
        <div
            class="stat rounded-2xl border border-base-200 bg-base-100 shadow-xl shadow-base-200/40 hover:shadow-md transition-all"
        >
            <div class="stat-figure text-primary opacity-30">
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    class="inline-block w-8 h-8 stroke-current"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"
                    />
                </svg>
            </div>
            <div class="stat-title text-xs font-bold uppercase tracking-wider">
                Total Kebijakan
            </div>
            <div class="stat-value text-primary text-2xl">
                {policies.length}
            </div>
        </div>

        <div
            class="stat rounded-2xl border border-base-200 bg-base-100 shadow-xl shadow-base-200/40 hover:shadow-md transition-all"
        >
            <div class="stat-figure text-secondary opacity-30">
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    class="inline-block w-8 h-8 stroke-current"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"
                    />
                </svg>
            </div>
            <div class="stat-title text-xs font-bold uppercase tracking-wider">
                Peran Aktif
            </div>
            <div class="stat-value text-secondary text-2xl">{roles.length}</div>
        </div>

        <div
            class="stat rounded-2xl border border-base-200 bg-base-100 shadow-xl shadow-base-200/40 hover:shadow-md transition-all"
        >
            <div class="stat-figure text-accent opacity-30">
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    class="inline-block w-8 h-8 stroke-current"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4m0 5c0 2.21-3.582 4-8 4s-8-1.79-8-4"
                    />
                </svg>
            </div>
            <div class="stat-title text-xs font-bold uppercase tracking-wider">
                Sumber Daya
            </div>
            <div class="stat-value text-accent text-2xl">
                {resources.length}
            </div>
        </div>
    </div>

    <div
        class="mb-8 rounded-2xl bg-info/10 border border-info/20 p-4 text-info flex items-start gap-3 shadow-sm shadow-info/5"
    >
        <svg
            xmlns="http://www.w3.org/2000/svg"
            class="h-6 w-6 shrink-0 mt-0.5"
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
        <div class="text-sm font-bold">
            <p>
                Klik kotak centang untuk memperbarui hak akses. Perubahan
                disimpan secara otomatis (real-time).
            </p>
        </div>
    </div>

    <AddRoleForm {isAdding} onSubmit={handleAddNewRoleWrapper} />

    <div
        class="card border border-base-200 bg-base-100 shadow-xl shadow-base-200/40 mt-6 overflow-hidden rounded-2xl"
    >
        <div class="overflow-x-auto">
            <table class="table w-full">
                <thead class="bg-base-200/50 text-base-content/70">
                    <tr class="border-b-2 border-base-200">
                        <th
                            class="font-bold text-xs uppercase tracking-widest py-5 pl-8 w-64"
                        >
                            Modul / Sumber Daya
                        </th>
                        {#each roles as role (role.id)}
                            <th
                                class="text-center py-5 min-w-50 border-l border-base-200/50"
                            >
                                <div
                                    class="flex items-center justify-center gap-3"
                                >
                                    <RoleBadge
                                        name={role.name}
                                        showText={true}
                                    />
                                    {#if auth.hasAllPermissions( ["role:delete"], )}
                                        <button
                                            type="button"
                                            class="btn btn-circle btn-ghost btn-xs text-error hover:bg-error/10 transition-colors"
                                            onclick={() =>
                                                confirmDeleteRole(role)}
                                            title="Hapus Peran"
                                        >
                                            <svg
                                                xmlns="http://www.w3.org/2000/svg"
                                                class="h-3.5 w-3.5"
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
                                        </button>
                                    {/if}
                                </div>
                            </th>
                        {/each}
                    </tr>
                </thead>
                <tbody class="divide-y divide-base-200">
                    {#each resources as resource}
                        <tr class="hover:bg-base-200/20 transition-colors">
                            <td
                                class="font-extrabold capitalize text-base-content align-top py-6 pl-8"
                            >
                                <div class="flex items-center gap-3">
                                    <div
                                        class="w-1.5 h-1.5 rounded-full bg-primary/40"
                                    ></div>
                                    {resource.name.replace("-", " ")}
                                </div>
                            </td>

                            {#each roles as role}
                                <td
                                    class="align-top py-6 border-l border-base-200/50"
                                >
                                    <div
                                        class="grid grid-cols-2 gap-y-3 gap-x-3 px-2"
                                    >
                                        {#each resource.actions as action}
                                            {@const isChecked = hasPolicy(
                                                role.name,
                                                resource.name,
                                                action,
                                            )}
                                            {@const isUpdating = isCellUpdating(
                                                role.name,
                                                resource.name,
                                                action,
                                            )}

                                            <label
                                                class="flex cursor-pointer items-center gap-2 text-[11px] p-2 rounded-xl transition-all duration-200 {isChecked
                                                    ? 'bg-primary/5'
                                                    : 'hover:bg-base-200'} {isUpdating
                                                    ? 'opacity-30 pointer-events-none'
                                                    : 'opacity-100'}"
                                            >
                                                <input
                                                    type="checkbox"
                                                    class="checkbox checkbox-xs checkbox-primary rounded shadow-sm"
                                                    checked={isChecked}
                                                    disabled={isUpdating ||
                                                        !auth.hasAllPermissions(
                                                            ["policy:update"],
                                                        )}
                                                    onchange={(e) =>
                                                        handleCheckboxChange(
                                                            role.name,
                                                            resource.name,
                                                            action,
                                                            e.currentTarget
                                                                .checked,
                                                        )}
                                                />
                                                <span
                                                    class="font-bold {isChecked
                                                        ? 'text-primary'
                                                        : 'text-base-content/50'} select-none whitespace-nowrap"
                                                >
                                                    {getActionLabel(action)}
                                                </span>
                                            </label>
                                        {/each}
                                    </div>
                                </td>
                            {/each}
                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>

        <div
            class="bg-base-200/40 p-5 border-t border-base-200 flex items-center justify-between flex-wrap gap-4"
        >
            <span
                class="text-xs font-bold uppercase tracking-widest text-base-content/40"
                >Ringkasan Peran</span
            >
            <div class="flex flex-wrap gap-2">
                {#each roles as role (role.id)}
                    <RoleBadge name={role.name} showText={true} />
                {/each}
            </div>
        </div>
    </div>
</div>
