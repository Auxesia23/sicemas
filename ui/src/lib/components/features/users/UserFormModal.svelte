<script lang="ts">
    import type { RoleResponse } from "$lib/schemas/user";

    interface Props {
        show: boolean;
        isEditMode: boolean;
        isLoading: boolean;
        formData: Record<string, any>;
        roles: RoleResponse[];
        modalError: string;
        errors: Record<string, string | undefined>;
        onSubmit: () => void;
        onCancel: () => void;
    }

    let {
        show = false,
        isEditMode = false,
        isLoading = false,
        formData = $bindable(),
        roles = [],
        modalError = "",
        errors = $bindable(),
        onSubmit,
        onCancel,
    }: Props = $props();

    function clearError(field: string) {
        if (errors[field]) {
            errors[field] = undefined;
        }
    }
</script>

<dialog class="modal modal-bottom sm:modal-middle" open={show}>
    <div
        class="modal-box p-0 max-w-2xl rounded-2xl shadow-2xl overflow-hidden bg-base-100 flex flex-col max-h-[90vh]"
    >
        <div class="bg-base-200/40 px-6 py-5 border-b border-base-200 sm:px-8">
            <div class="flex items-center gap-4">
                <div
                    class="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-linear-to-br from-primary/20 to-secondary/20 text-primary shadow-sm"
                >
                    {#if isEditMode}
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
                                d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                            />
                        </svg>
                    {:else}
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
                                d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z"
                            />
                        </svg>
                    {/if}
                </div>
                <div>
                    <h3 class="text-xl font-extrabold text-base-content">
                        {isEditMode
                            ? "Edit Data Petugas"
                            : "Registrasi Petugas Baru"}
                    </h3>
                    <p class="text-sm font-medium text-base-content/60 mt-0.5">
                        Lengkapi informasi di bawah ini dengan benar.
                    </p>
                </div>
            </div>
        </div>

        <div class="px-6 py-6 sm:px-8 overflow-y-auto flex-1">
            {#if modalError}
                <div
                    class="mb-6 alert py-3 text-sm bg-error/10 text-error border border-error/20 rounded-xl"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="h-5 w-5 shrink-0 stroke-current"
                        fill="none"
                        viewBox="0 0 24 24"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                        />
                    </svg>
                    <span class="font-medium">{modalError}</span>
                </div>
            {/if}

            <div class="grid grid-cols-1 gap-5 md:grid-cols-2">
                <div class="form-control">
                    <label class="label pb-1.5 pt-0" for="nip">
                        <span class="label-text font-bold text-base-content/80">
                            NIP / NIK <span class="text-error">*</span>
                        </span>
                    </label>
                    <input
                        id="nip"
                        type="text"
                        class="input input-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors {errors.nip
                            ? 'input-error focus:border-error'
                            : 'focus:border-primary'}"
                        bind:value={formData.nip}
                        placeholder="16 atau 18 digit angka"
                        maxlength="18"
                        oninput={() => clearError("nip")}
                        disabled={isLoading}
                    />
                    {#if errors.nip}
                        <span class="mt-1.5 text-xs font-medium text-error"
                            >{errors.nip}</span
                        >
                    {/if}
                </div>

                <div class="form-control">
                    <label class="label pb-1.5 pt-0" for="nama">
                        <span class="label-text font-bold text-base-content/80">
                            Nama Lengkap <span class="text-error">*</span>
                        </span>
                    </label>
                    <input
                        id="nama"
                        type="text"
                        class="input input-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors {errors.nama_lengkap
                            ? 'input-error focus:border-error'
                            : 'focus:border-primary'}"
                        bind:value={formData.nama_lengkap}
                        placeholder="Nama lengkap beserta gelar"
                        oninput={() => clearError("nama_lengkap")}
                        disabled={isLoading}
                    />
                    {#if errors.nama_lengkap}
                        <span class="mt-1.5 text-xs font-medium text-error"
                            >{errors.nama_lengkap}</span
                        >
                    {/if}
                </div>

                <div class="form-control">
                    <label class="label pb-1.5 pt-0" for="jabatan">
                        <span class="label-text font-bold text-base-content/80">
                            Jabatan <span class="text-error">*</span>
                        </span>
                    </label>
                    <input
                        id="jabatan"
                        type="text"
                        class="input input-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors {errors.jabatan
                            ? 'input-error focus:border-error'
                            : 'focus:border-primary'}"
                        bind:value={formData.jabatan}
                        placeholder="Contoh: Penghulu"
                        oninput={() => clearError("jabatan")}
                        disabled={isLoading}
                    />
                    {#if errors.jabatan}
                        <span class="mt-1.5 text-xs font-medium text-error"
                            >{errors.jabatan}</span
                        >
                    {/if}
                </div>

                <div class="form-control">
                    <label class="label pb-1.5 pt-0" for="unit_kerja">
                        <span class="label-text font-bold text-base-content/80">
                            Unit Kerja <span
                                class="text-xs font-normal opacity-50"
                                >(Opsional)</span
                            >
                        </span>
                    </label>
                    <input
                        id="unit_kerja"
                        type="text"
                        class="input input-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors {errors.unit_kerja
                            ? 'input-error focus:border-error'
                            : 'focus:border-primary'}"
                        bind:value={formData.unit_kerja}
                        placeholder="Default: KUA Ciemas"
                        oninput={() => clearError("unit_kerja")}
                        disabled={isLoading}
                    />
                    {#if errors.unit_kerja}
                        <span class="mt-1.5 text-xs font-medium text-error"
                            >{errors.unit_kerja}</span
                        >
                    {/if}
                </div>

                <div class="form-control">
                    <label class="label pb-1.5 pt-0" for="peran">
                        <span class="label-text font-bold text-base-content/80">
                            Peran <span class="text-error">*</span>
                        </span>
                    </label>
                    <select
                        id="peran"
                        class="select select-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors {errors.peran
                            ? 'select-error focus:border-error'
                            : 'focus:border-primary'}"
                        bind:value={formData.peran}
                        onchange={() => clearError("peran")}
                        disabled={isLoading}
                    >
                        <option value="" disabled selected={!formData.peran}
                            >Pilih peran sistem</option
                        >
                        {#each roles as role (role.id)}
                            <option value={role.name}>{role.name}</option>
                        {/each}
                    </select>
                    {#if errors.peran}
                        <span class="mt-1.5 text-xs font-medium text-error"
                            >{errors.peran}</span
                        >
                    {/if}
                </div>

                <div class="form-control">
                    <label class="label pb-1.5 pt-0" for="email">
                        <span class="label-text font-bold text-base-content/80">
                            Email <span class="text-xs font-normal opacity-50"
                                >(Opsional)</span
                            >
                        </span>
                    </label>
                    <input
                        id="email"
                        type="email"
                        class="input input-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors {errors.email
                            ? 'input-error focus:border-error'
                            : 'focus:border-primary'}"
                        bind:value={formData.email}
                        placeholder="petugas@kemenag.go.id"
                        oninput={() => clearError("email")}
                        disabled={isLoading}
                    />
                    {#if errors.email}
                        <span class="mt-1.5 text-xs font-medium text-error"
                            >{errors.email}</span
                        >
                    {/if}
                </div>

                <div class="form-control md:col-span-2">
                    <label class="label pb-1.5 pt-0" for="telp">
                        <span class="label-text font-bold text-base-content/80">
                            Nomor Telepon <span class="text-error">*</span>
                        </span>
                    </label>
                    <input
                        id="telp"
                        type="tel"
                        class="input input-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors {errors.nomor_telepon
                            ? 'input-error focus:border-error'
                            : 'focus:border-primary'}"
                        bind:value={formData.nomor_telepon}
                        placeholder="08xxxxxxxxxx"
                        oninput={() => clearError("nomor_telepon")}
                        disabled={isLoading}
                    />
                    {#if errors.nomor_telepon}
                        <span class="mt-1.5 text-xs font-medium text-error"
                            >{errors.nomor_telepon}</span
                        >
                    {/if}
                </div>
            </div>
        </div>

        <div
            class="bg-base-200/40 px-6 py-5 border-t border-base-200 sm:px-8 flex flex-col-reverse sm:flex-row justify-end gap-3 sticky bottom-0"
        >
            <button
                type="button"
                class="btn btn-ghost rounded-xl hover:bg-base-300 w-full sm:w-auto"
                onclick={onCancel}
                disabled={isLoading}
            >
                Batal
            </button>
            <button
                type="button"
                class="btn btn-primary rounded-xl px-8 shadow-lg shadow-primary/20 transition-transform hover:-translate-y-0.5 w-full sm:w-auto text-primary-content"
                onclick={onSubmit}
                disabled={isLoading}
            >
                {#if isLoading}
                    <span class="loading loading-sm loading-spinner"></span>
                    Menyimpan...
                {:else}
                    {isEditMode ? "Perbarui Data" : "Simpan Data"}
                {/if}
            </button>
        </div>
    </div>

    <button
        type="button"
        aria-label="Close modal"
        class="modal-backdrop bg-base-content/20 backdrop-blur-sm"
        onclick={onCancel}
        disabled={isLoading}
    ></button>
</dialog>
