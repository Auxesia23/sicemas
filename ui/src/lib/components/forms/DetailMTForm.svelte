<script lang="ts">
    import { DEFAULT_DETAIL_MT } from "$lib/schemas/siteInput";
    import {
        addToArray,
        removeFromArray,
        handleEnterKey,
        clearNestedError,
    } from "$lib/utils/formHelpers";

    interface Props {
        detail?: Record<string, any>;
        errors?: Record<string, any>;
    }

    let { detail = $bindable({}), errors = $bindable({}) }: Props = $props();

    let newSekretaris = $state("");
    let newBendahara = $state("");
    let newPenceramahNama = $state("");
    let newPenceramahMateri = $state("");

    if (!detail?.pengurus) {
        detail = structuredClone(DEFAULT_DETAIL_MT);
    }

    function addSekretaris() {
        newSekretaris = addToArray(detail.pengurus.sekretaris, newSekretaris);
        clearNestedError(errors, "pengurus", "sekretaris");
    }

    function addBendahara() {
        newBendahara = addToArray(detail.pengurus.bendahara, newBendahara);
        clearNestedError(errors, "pengurus", "bendahara");
    }

    function addPenceramah() {
        if (newPenceramahNama.trim() && newPenceramahMateri.trim()) {
            if (!detail.penceramah) detail.penceramah = [];
            detail.penceramah.push({
                nama: newPenceramahNama.trim(),
                materi: newPenceramahMateri.trim(),
            });
            newPenceramahNama = "";
            newPenceramahMateri = "";
            clearNestedError(errors, "penceramah");
        }
    }

    function removePenceramah(index: number) {
        detail.penceramah.splice(index, 1);
    }
</script>

<div
    class="card rounded-2xl border border-base-200 bg-base-100 shadow-xl shadow-base-200/40 overflow-hidden"
>
    <div class="card-body p-5 sm:p-6">
        <div class="flex items-center gap-3 mb-5">
            <div
                class="flex h-10 w-10 items-center justify-center rounded-xl bg-linear-to-br from-primary/20 to-secondary/20 text-primary font-bold shadow-inner"
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
                        d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
                    />
                </svg>
            </div>
            <h2
                class="card-title text-xl font-extrabold tracking-tight text-base-content"
            >
                Pengurus <span class="text-error">*</span>
            </h2>
        </div>

        <div class="mb-5">
            <label class="label py-2" for="pengurus_ketua">
                <span
                    class="label-text font-bold text-base-content/80 text-xs uppercase tracking-wider"
                >
                    Ketua <span class="text-error">*</span>
                </span>
            </label>
            <input
                id="pengurus_ketua"
                type="text"
                class="input input-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-11 {errors
                    ?.pengurus?.ketua
                    ? 'input-error'
                    : ''}"
                bind:value={detail.pengurus.ketua}
                oninput={() => clearNestedError(errors, "pengurus", "ketua")}
            />
            {#if errors?.pengurus?.ketua}
                <span class="mt-1 text-xs text-error font-medium"
                    >{errors.pengurus.ketua}</span
                >
            {/if}
        </div>

        <div class="grid grid-cols-1 gap-5 sm:grid-cols-2">
            <div class="form-control">
                <label class="label py-2" for="sekretaris_input">
                    <span
                        class="label-text font-bold text-base-content/80 text-xs uppercase tracking-wider"
                    >
                        Sekretaris <span class="text-error">*</span>
                    </span>
                </label>
                <div class="flex gap-2">
                    <input
                        id="sekretaris_input"
                        type="text"
                        class="input input-bordered flex-1 rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-11 {errors
                            ?.pengurus?.sekretaris
                            ? 'input-error'
                            : ''}"
                        placeholder="Nama Sekretaris"
                        bind:value={newSekretaris}
                        onkeydown={(e) => handleEnterKey(e, addSekretaris)}
                    />
                    <button
                        type="button"
                        class="btn btn-primary rounded-xl px-5 shadow-lg shadow-primary/20 transition-transform hover:-translate-y-0.5 min-h-11"
                        onclick={addSekretaris}
                    >
                        Tambah
                    </button>
                </div>
                <div class="mt-3 flex flex-wrap gap-2">
                    {#each detail.pengurus?.sekretaris || [] as person, index (index)}
                        <div
                            class="badge font-medium px-3 py-3 rounded-xl bg-info/10 text-info border-info/20 gap-2"
                        >
                            {person}
                            <button
                                type="button"
                                class="btn btn-circle btn-ghost btn-xs h-5 w-5 min-h-0 p-0 hover:bg-info/20"
                                aria-label="Hapus sekretaris {person}"
                                onclick={() =>
                                    removeFromArray(
                                        detail.pengurus.sekretaris,
                                        index,
                                    )}
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
                                        d="M6 18L18 6M6 6l12 12"
                                    />
                                </svg>
                            </button>
                        </div>
                    {/each}
                </div>
                {#if errors?.pengurus?.sekretaris}
                    <span class="mt-1 text-xs text-error font-medium"
                        >{errors.pengurus.sekretaris}</span
                    >
                {/if}
            </div>

            <div class="form-control">
                <label class="label py-2" for="bendahara_input">
                    <span
                        class="label-text font-bold text-base-content/80 text-xs uppercase tracking-wider"
                    >
                        Bendahara <span class="text-error">*</span>
                    </span>
                </label>
                <div class="flex gap-2">
                    <input
                        id="bendahara_input"
                        type="text"
                        class="input input-bordered flex-1 rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-11 {errors
                            ?.pengurus?.bendahara
                            ? 'input-error'
                            : ''}"
                        placeholder="Nama Bendahara"
                        bind:value={newBendahara}
                        onkeydown={(e) => handleEnterKey(e, addBendahara)}
                    />
                    <button
                        type="button"
                        class="btn btn-primary rounded-xl px-5 shadow-lg shadow-primary/20 transition-transform hover:-translate-y-0.5 min-h-11"
                        onclick={addBendahara}
                    >
                        Tambah
                    </button>
                </div>
                <div class="mt-3 flex flex-wrap gap-2">
                    {#each detail.pengurus?.bendahara || [] as person, index (index)}
                        <div
                            class="badge font-medium px-3 py-3 rounded-xl bg-warning/10 text-warning border-warning/20 gap-2"
                        >
                            {person}
                            <button
                                type="button"
                                class="btn btn-circle btn-ghost btn-xs h-5 w-5 min-h-0 p-0 hover:bg-warning/20"
                                aria-label="Hapus bendahara {person}"
                                onclick={() =>
                                    removeFromArray(
                                        detail.pengurus.bendahara,
                                        index,
                                    )}
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
                                        d="M6 18L18 6M6 6l12 12"
                                    />
                                </svg>
                            </button>
                        </div>
                    {/each}
                </div>
                {#if errors?.pengurus?.bendahara}
                    <span class="mt-1 text-xs text-error font-medium"
                        >{errors.pengurus.bendahara}</span
                    >
                {/if}
            </div>
        </div>
    </div>
</div>

{#if typeof errors?.kehadiran === "string"}
    <div class="px-2">
        <span
            class="text-sm font-medium text-error bg-error/10 rounded-xl p-3 border border-error/20 block"
            >{errors.kehadiran}</span
        >
    </div>
{/if}

<div
    class="card rounded-2xl border border-base-200 bg-base-100 shadow-xl shadow-base-200/40 overflow-hidden {errors?.kehadiran
        ? 'border-error border-2'
        : ''}"
>
    <div class="card-body p-5 sm:p-6">
        <div class="flex items-center gap-3 mb-5">
            <div
                class="flex h-10 w-10 items-center justify-center rounded-xl bg-linear-to-br from-info/20 to-primary/20 text-info font-bold shadow-inner"
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
                        d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
                    />
                </svg>
            </div>
            <h2
                class="card-title text-xl font-extrabold tracking-tight text-base-content"
            >
                Kehadiran Pria
            </h2>
        </div>

        <div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
            <div class="form-control">
                <label class="label py-2" for="pria_dewasa">
                    <span
                        class="label-text font-bold text-base-content/80 text-xs uppercase tracking-wider"
                        >Jml. Dewasa</span
                    >
                </label>
                <input
                    id="pria_dewasa"
                    type="number"
                    min="0"
                    class="input input-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-11"
                    placeholder="0"
                    bind:value={detail.kehadiran_pria.jumlah_dewasa}
                />
            </div>
            <div class="form-control">
                <label class="label py-2" for="pria_remaja">
                    <span
                        class="label-text font-bold text-base-content/80 text-xs uppercase tracking-wider"
                        >Jml. Remaja</span
                    >
                </label>
                <input
                    id="pria_remaja"
                    type="number"
                    min="0"
                    class="input input-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-11"
                    placeholder="0"
                    bind:value={detail.kehadiran_pria.jumlah_remaja}
                />
            </div>
            <div class="form-control">
                <label class="label py-2" for="pria_waktu">
                    <span
                        class="label-text font-bold text-base-content/80 text-xs uppercase tracking-wider"
                        >Waktu Pengajian</span
                    >
                </label>
                <select
                    id="pria_waktu"
                    class="select select-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-11"
                    bind:value={detail.kehadiran_pria.waktu_pengajian}
                >
                    <option value="">Pilih...</option>
                    <option value="Pagi">Pagi</option>
                    <option value="Sore">Sore</option>
                    <option value="Malam">Malam</option>
                </select>
            </div>
        </div>
    </div>
</div>

<div
    class="card rounded-2xl border border-base-200 bg-base-100 shadow-xl shadow-base-200/40 overflow-hidden {errors?.kehadiran
        ? 'border-error border-2'
        : ''}"
>
    <div class="card-body p-5 sm:p-6">
        <div class="flex items-center gap-3 mb-5">
            <div
                class="flex h-10 w-10 items-center justify-center rounded-xl bg-linear-to-br from-secondary/20 to-accent/20 text-secondary font-bold shadow-inner"
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
                        d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"
                    />
                </svg>
            </div>
            <h2
                class="card-title text-xl font-extrabold tracking-tight text-base-content"
            >
                Kehadiran Wanita
            </h2>
        </div>

        <div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
            <div class="form-control">
                <label class="label py-2" for="wanita_dewasa">
                    <span
                        class="label-text font-bold text-base-content/80 text-xs uppercase tracking-wider"
                        >Jml. Dewasa</span
                    >
                </label>
                <input
                    id="wanita_dewasa"
                    type="number"
                    min="0"
                    class="input input-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-11"
                    bind:value={detail.kehadiran_wanita.jumlah_dewasa}
                />
            </div>
            <div class="form-control">
                <label class="label py-2" for="wanita_remaja">
                    <span
                        class="label-text font-bold text-base-content/80 text-xs uppercase tracking-wider"
                        >Jml. Remaja</span
                    >
                </label>
                <input
                    id="wanita_remaja"
                    type="number"
                    min="0"
                    class="input input-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-11"
                    bind:value={detail.kehadiran_wanita.jumlah_remaja}
                />
            </div>
            <div class="form-control">
                <label class="label py-2" for="wanita_waktu">
                    <span
                        class="label-text font-bold text-base-content/80 text-xs uppercase tracking-wider"
                        >Waktu Pengajian</span
                    >
                </label>
                <select
                    id="wanita_waktu"
                    class="select select-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-11"
                    bind:value={detail.kehadiran_wanita.waktu_pengajian}
                >
                    <option value="">Pilih...</option>
                    <option value="Pagi">Pagi</option>
                    <option value="Sore">Sore</option>
                    <option value="Malam">Malam</option>
                </select>
            </div>
        </div>
    </div>
</div>

<div
    class="card rounded-2xl border border-base-200 bg-base-100 shadow-xl shadow-base-200/40 overflow-hidden"
>
    <div class="card-body p-5 sm:p-6">
        <div class="flex items-center gap-3 mb-5">
            <div
                class="flex h-10 w-10 items-center justify-center rounded-xl bg-linear-to-br from-warning/20 to-error/20 text-warning font-bold shadow-inner"
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
                        d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z"
                    />
                </svg>
            </div>
            <h2
                class="card-title text-xl font-extrabold tracking-tight text-base-content"
            >
                Daftar Penceramah <span class="text-error">*</span>
            </h2>
        </div>

        <div>
            <label class="label py-2" for="penceramah_inputs">
                <span
                    class="label-text font-bold text-base-content/80 text-xs uppercase tracking-wider"
                    >Tambah Penceramah</span
                >
            </label>
            <div class="mb-3 grid grid-cols-1 gap-3 sm:grid-cols-12">
                <div class="form-control sm:col-span-5">
                    <input
                        id="penceramah_nama"
                        type="text"
                        class="input input-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-11 {errors?.penceramah
                            ? 'input-error'
                            : ''}"
                        placeholder="Nama penceramah..."
                        bind:value={newPenceramahNama}
                        onkeydown={(e) => handleEnterKey(e, addPenceramah)}
                        oninput={() => clearNestedError(errors, "penceramah")}
                    />
                </div>
                <div class="form-control sm:col-span-5">
                    <input
                        id="penceramah_materi"
                        type="text"
                        class="input input-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-11 {errors?.penceramah
                            ? 'input-error'
                            : ''}"
                        placeholder="Materi..."
                        bind:value={newPenceramahMateri}
                        onkeydown={(e) => handleEnterKey(e, addPenceramah)}
                        oninput={() => clearNestedError(errors, "penceramah")}
                    />
                </div>
                <button
                    type="button"
                    class="btn btn-primary rounded-xl px-5 shadow-lg shadow-primary/20 transition-transform hover:-translate-y-0.5 min-h-11 sm:col-span-2"
                    onclick={addPenceramah}
                >
                    Tambah
                </button>
            </div>
            {#if errors?.penceramah && typeof errors.penceramah === "string"}
                <span class="mt-1 block text-xs text-error font-medium"
                    >{errors.penceramah}</span
                >
            {/if}

            <div class="mt-4 flex flex-col gap-3">
                {#each detail.penceramah || [] as penceramah, index (index)}
                    <div
                        class="flex items-center justify-between rounded-xl bg-base-200/30 border border-base-200/50 p-4 hover:bg-base-200/50 transition-colors"
                    >
                        <div class="flex-1">
                            <div class="font-semibold text-base-content">
                                {penceramah.nama}
                            </div>
                            <div
                                class="text-sm font-medium text-base-content/60 flex items-center gap-1.5 mt-0.5"
                            >
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    class="h-4 w-4 opacity-50"
                                    fill="none"
                                    viewBox="0 0 24 24"
                                    stroke="currentColor"
                                >
                                    <path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"
                                    />
                                </svg>
                                {penceramah.materi}
                            </div>
                        </div>
                        <button
                            type="button"
                            class="btn btn-circle btn-xs btn-error shadow-lg"
                            aria-label="Hapus penceramah {penceramah.nama}"
                            onclick={() => removePenceramah(index)}
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
                                    d="M6 18L18 6M6 6l12 12"
                                />
                            </svg>
                        </button>
                    </div>
                {/each}
                {#if detail.penceramah?.length === 0}
                    <div class="text-center py-6">
                        <div
                            class="inline-flex items-center justify-center w-12 h-12 rounded-full bg-base-200/50 mb-2"
                        >
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                class="h-6 w-6 text-base-content/30"
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
                        <p class="text-sm font-medium text-base-content/50">
                            Belum ada penceramah yang ditambahkan
                        </p>
                    </div>
                {/if}
            </div>
        </div>
    </div>
</div>
