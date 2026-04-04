<script lang="ts">
    import { DEFAULT_DETAIL_MASJID } from "$lib/schemas/siteInput";
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

    let newJenisBuku = $state("");
    let newNamaImam = $state("");
    let newNamaMuazdin = $state("");
    let newSekretaris = $state("");
    let newBendahara = $state("");
    let newFasilitasUmum = $state("");
    let newFasilitasRamanAnak = $state("");
    let newFasilitasDisabilitas = $state("");
    let newKegiatan = $state("");

    if (!detail?.perpustakaan) {
        detail = structuredClone(DEFAULT_DETAIL_MASJID);
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
                SDM Masjid <span class="text-error">*</span>
            </h2>
        </div>

        {#if errors?.sdm_masjid && typeof errors.sdm_masjid === "string"}
            <div
                class="mb-4 text-sm font-medium text-error bg-error/10 rounded-xl p-3 border border-error/20"
            >
                {errors.sdm_masjid}
            </div>
        {/if}

        <div class="grid grid-cols-2 gap-3 sm:grid-cols-3 sm:gap-4">
            {#each [{ id: "pengurus", label: "Jml Pengurus", path: "jumlah_pengurus" }, { id: "imam", label: "Jml Imam", path: "jumlah_imam" }, { id: "khotib", label: "Jml Khotib", path: "jumlah_khotib" }, { id: "muadzin", label: "Jml Muadzin", path: "jumlah_muadzin" }, { id: "remaja", label: "Jml Remaja", path: "jumlah_remaja" }, { id: "jemaah", label: "Jml Jemaah", path: "jumlah_jemaah" }] as item}
                <div class="form-control">
                    <label class="label py-2" for="sdm_{item.id}">
                        <span
                            class="label-text font-bold text-base-content/80 text-xs uppercase tracking-wider"
                        >
                            {item.label} <span class="text-error">*</span>
                        </span>
                    </label>
                    <input
                        id="sdm_{item.id}"
                        type="number"
                        min="0"
                        class="input input-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-11 {errors
                            ?.sdm_masjid?.[item.path]
                            ? 'input-error'
                            : ''}"
                        bind:value={detail.sdm_masjid[item.path]}
                        oninput={() =>
                            clearNestedError(errors, "sdm_masjid", item.path)}
                    />
                    {#if errors?.sdm_masjid?.[item.path]}
                        <span class="mt-1 text-xs text-error font-medium"
                            >{errors.sdm_masjid[item.path]}</span
                        >
                    {/if}
                </div>
            {/each}
        </div>
    </div>
</div>

<div
    class="card rounded-2xl border border-base-200 bg-base-100 shadow-xl shadow-base-200/40 overflow-hidden"
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
                Pengurus DKM <span class="text-error">*</span>
            </h2>
        </div>

        <div class="form-control mb-5">
            <label class="label py-2" for="pengurus_ketua">
                <span
                    class="label-text font-bold text-base-content/80 text-xs uppercase tracking-wider"
                >
                    Nama Ketua DKM <span class="text-error">*</span>
                </span>
            </label>
            <input
                id="pengurus_ketua"
                type="text"
                class="input input-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-11 {errors
                    ?.pengurus_dkm?.ketua
                    ? 'input-error'
                    : ''}"
                bind:value={detail.pengurus_dkm.ketua}
                oninput={() =>
                    clearNestedError(errors, "pengurus_dkm", "ketua")}
            />
            {#if errors?.pengurus_dkm?.ketua}
                <span class="mt-1 text-xs text-error font-medium"
                    >{errors.pengurus_dkm.ketua}</span
                >
            {/if}
        </div>

        <div class="grid grid-cols-1 gap-5 sm:grid-cols-2">
            <div class="form-control">
                <label class="label py-2" for="sek_in">
                    <span
                        class="label-text font-bold text-base-content/80 text-xs uppercase tracking-wider"
                    >
                        Sekretaris <span class="text-error">*</span>
                    </span>
                </label>
                <div class="flex gap-2">
                    <input
                        id="sek_in"
                        type="text"
                        class="input input-bordered flex-1 rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-11 {errors
                            ?.pengurus_dkm?.sekretaris
                            ? 'input-error'
                            : ''}"
                        placeholder="Tambah Sekretaris..."
                        bind:value={newSekretaris}
                        onkeydown={(e) =>
                            handleEnterKey(e, () => {
                                newSekretaris = addToArray(
                                    detail.pengurus_dkm.sekretaris,
                                    newSekretaris,
                                );
                                clearNestedError(
                                    errors,
                                    "pengurus_dkm",
                                    "sekretaris",
                                );
                            })}
                    />
                    <button
                        type="button"
                        class="btn btn-primary rounded-xl px-5 shadow-lg shadow-primary/20 transition-transform hover:-translate-y-0.5 min-h-11"
                        onclick={() => {
                            newSekretaris = addToArray(
                                detail.pengurus_dkm.sekretaris,
                                newSekretaris,
                            );
                            clearNestedError(
                                errors,
                                "pengurus_dkm",
                                "sekretaris",
                            );
                        }}
                    >
                        Tambah
                    </button>
                </div>
                <div class="mt-3 flex flex-wrap gap-2">
                    {#each detail.pengurus_dkm.sekretaris as p, i}
                        <div
                            class="badge font-medium px-3 py-3 rounded-xl bg-info/10 text-info border-info/20 gap-2"
                        >
                            {p}
                            <button
                                title="Hapus Sekretaris"
                                type="button"
                                class="btn btn-circle btn-ghost btn-xs h-5 w-5 min-h-0 p-0 hover:bg-info/20"
                                onclick={() =>
                                    removeFromArray(
                                        detail.pengurus_dkm.sekretaris,
                                        i,
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
                {#if errors?.pengurus_dkm?.sekretaris}
                    <span class="mt-1 text-xs text-error font-medium"
                        >{errors.pengurus_dkm.sekretaris}</span
                    >
                {/if}
            </div>

            <div class="form-control">
                <label class="label py-2" for="ben_in">
                    <span
                        class="label-text font-bold text-base-content/80 text-xs uppercase tracking-wider"
                    >
                        Bendahara <span class="text-error">*</span>
                    </span>
                </label>
                <div class="flex gap-2">
                    <input
                        id="ben_in"
                        type="text"
                        class="input input-bordered flex-1 rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-11 {errors
                            ?.pengurus_dkm?.bendahara
                            ? 'input-error'
                            : ''}"
                        placeholder="Tambah Bendahara..."
                        bind:value={newBendahara}
                        onkeydown={(e) =>
                            handleEnterKey(e, () => {
                                newBendahara = addToArray(
                                    detail.pengurus_dkm.bendahara,
                                    newBendahara,
                                );
                                clearNestedError(
                                    errors,
                                    "pengurus_dkm",
                                    "bendahara",
                                );
                            })}
                    />
                    <button
                        type="button"
                        class="btn btn-primary rounded-xl px-5 shadow-lg shadow-primary/20 transition-transform hover:-translate-y-0.5 min-h-11"
                        onclick={() => {
                            newBendahara = addToArray(
                                detail.pengurus_dkm.bendahara,
                                newBendahara,
                            );
                            clearNestedError(
                                errors,
                                "pengurus_dkm",
                                "bendahara",
                            );
                        }}
                    >
                        Tambah
                    </button>
                </div>
                <div class="mt-3 flex flex-wrap gap-2">
                    {#each detail.pengurus_dkm.bendahara as p, i}
                        <div
                            class="badge font-medium px-3 py-3 rounded-xl bg-warning/10 text-warning border-warning/20 gap-2"
                        >
                            {p}
                            <button
                                title="Hapus Bendahara"
                                type="button"
                                class="btn btn-circle btn-ghost btn-xs h-5 w-5 min-h-0 p-0 hover:bg-warning/20"
                                onclick={() =>
                                    removeFromArray(
                                        detail.pengurus_dkm.bendahara,
                                        i,
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
                {#if errors?.pengurus_dkm?.bendahara}
                    <span class="mt-1 text-xs text-error font-medium"
                        >{errors.pengurus_dkm.bendahara}</span
                    >
                {/if}
            </div>
        </div>

        <div class="mt-5 grid grid-cols-1 gap-5 sm:grid-cols-2">
            <div class="form-control">
                <label class="label py-2" for="imam_in">
                    <span
                        class="label-text font-bold text-base-content/80 text-xs uppercase tracking-wider"
                    >
                        Nama Imam <span class="text-error">*</span>
                    </span>
                </label>
                <div class="flex gap-2">
                    <input
                        id="imam_in"
                        type="text"
                        class="input input-bordered flex-1 rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-11 {errors
                            ?.pengurus_dkm?.nama_imam
                            ? 'input-error'
                            : ''}"
                        bind:value={newNamaImam}
                        onkeydown={(e) =>
                            handleEnterKey(e, () => {
                                newNamaImam = addToArray(
                                    detail.pengurus_dkm.nama_imam,
                                    newNamaImam,
                                );
                                clearNestedError(
                                    errors,
                                    "pengurus_dkm",
                                    "nama_imam",
                                );
                            })}
                    />
                    <button
                        type="button"
                        class="btn btn-primary rounded-xl px-5 shadow-lg shadow-primary/20 transition-transform hover:-translate-y-0.5 min-h-11"
                        onclick={() => {
                            newNamaImam = addToArray(
                                detail.pengurus_dkm.nama_imam,
                                newNamaImam,
                            );
                            clearNestedError(
                                errors,
                                "pengurus_dkm",
                                "nama_imam",
                            );
                        }}
                    >
                        Tambah
                    </button>
                </div>
                <div class="mt-3 flex flex-wrap gap-2">
                    {#each detail.pengurus_dkm.nama_imam as p, i}
                        <div
                            class="badge font-medium px-3 py-3 rounded-xl bg-primary/10 text-primary border-primary/20 gap-2"
                        >
                            {p}
                            <button
                                title="Hapus Nama Imam"
                                type="button"
                                class="btn btn-circle btn-ghost btn-xs h-5 w-5 min-h-0 p-0 hover:bg-primary/20"
                                onclick={() =>
                                    removeFromArray(
                                        detail.pengurus_dkm.nama_imam,
                                        i,
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
                {#if errors?.pengurus_dkm?.nama_imam}
                    <span class="mt-1 text-xs text-error font-medium"
                        >{errors.pengurus_dkm.nama_imam}</span
                    >
                {/if}
            </div>

            <div class="form-control">
                <label class="label py-2" for="mua_in">
                    <span
                        class="label-text font-bold text-base-content/80 text-xs uppercase tracking-wider"
                    >
                        Nama Muadzin <span class="text-error">*</span>
                    </span>
                </label>
                <div class="flex gap-2">
                    <input
                        id="mua_in"
                        type="text"
                        class="input input-bordered flex-1 rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-11 {errors
                            ?.pengurus_dkm?.nama_muazdin
                            ? 'input-error'
                            : ''}"
                        bind:value={newNamaMuazdin}
                        onkeydown={(e) =>
                            handleEnterKey(e, () => {
                                newNamaMuazdin = addToArray(
                                    detail.pengurus_dkm.nama_muazdin,
                                    newNamaMuazdin,
                                );
                                clearNestedError(
                                    errors,
                                    "pengurus_dkm",
                                    "nama_muazdin",
                                );
                            })}
                    />
                    <button
                        type="button"
                        class="btn btn-primary rounded-xl px-5 shadow-lg shadow-primary/20 transition-transform hover:-translate-y-0.5 min-h-11"
                        onclick={() => {
                            newNamaMuazdin = addToArray(
                                detail.pengurus_dkm.nama_muazdin,
                                newNamaMuazdin,
                            );
                            clearNestedError(
                                errors,
                                "pengurus_dkm",
                                "nama_muazdin",
                            );
                        }}
                    >
                        Tambah
                    </button>
                </div>
                <div class="mt-3 flex flex-wrap gap-2">
                    {#each detail.pengurus_dkm.nama_muazdin as p, i}
                        <div
                            class="badge font-medium px-3 py-3 rounded-xl bg-primary/10 text-primary border-primary/20 gap-2"
                        >
                            {p}
                            <button
                                title="Hapus Nama Muazdin"
                                type="button"
                                class="btn btn-circle btn-ghost btn-xs h-5 w-5 min-h-0 p-0 hover:bg-primary/20"
                                onclick={() =>
                                    removeFromArray(
                                        detail.pengurus_dkm.nama_muazdin,
                                        i,
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
                {#if errors?.pengurus_dkm?.nama_muazdin}
                    <span class="mt-1 text-xs text-error font-medium"
                        >{errors.pengurus_dkm.nama_muazdin}</span
                    >
                {/if}
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
                        d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"
                    />
                </svg>
            </div>
            <h2
                class="card-title text-xl font-extrabold tracking-tight text-base-content"
            >
                Perpustakaan
            </h2>
        </div>

        <div class="grid grid-cols-2 gap-3 sm:grid-cols-3 sm:gap-4">
            <div class="form-control">
                <label class="label py-2" for="perpustakaan_luas_m2">
                    <span
                        class="label-text font-bold text-base-content/80 text-xs uppercase tracking-wider"
                        >Luas (m²)</span
                    >
                </label>
                <input
                    id="perpustakaan_luas_m2"
                    type="number"
                    min="0"
                    class="input input-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-11"
                    bind:value={detail.perpustakaan.luas_m2}
                />
            </div>
            <div class="form-control">
                <label class="label py-2" for="perpustakaan_jumlah_pengurus">
                    <span
                        class="label-text font-bold text-base-content/80 text-xs uppercase tracking-wider"
                        >Jml. Pengurus</span
                    >
                </label>
                <input
                    id="perpustakaan_jumlah_pengurus"
                    type="number"
                    min="0"
                    class="input input-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-11"
                    bind:value={detail.perpustakaan.jumlah_pengurus}
                />
            </div>
            <div class="form-control col-span-2 sm:col-span-1">
                <label class="label py-2" for="perpustakaan_kondisi">
                    <span
                        class="label-text font-bold text-base-content/80 text-xs uppercase tracking-wider"
                        >Kondisi</span
                    >
                </label>
                <select
                    id="perpustakaan_kondisi"
                    class="select select-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-11"
                    bind:value={detail.perpustakaan.kondisi}
                >
                    <option value="">Pilih...</option>
                    <option value="Baik">Baik</option>
                    <option value="Cukup Baik">Cukup Baik</option>
                    <option value="Kurang Baik">Kurang Baik</option>
                    <option value="Rusak">Rusak</option>
                </select>
            </div>
        </div>

        <div class="mt-5">
            <label class="label py-2" for="jenis_buku_input">
                <span
                    class="label-text font-bold text-base-content/80 text-xs uppercase tracking-wider"
                    >Jenis Buku</span
                >
            </label>
            <div class="mb-3 flex gap-2">
                <input
                    id="jenis_buku_input"
                    type="text"
                    class="input input-bordered flex-1 rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-11"
                    placeholder="Tambah jenis buku..."
                    bind:value={newJenisBuku}
                    onkeydown={(e) =>
                        handleEnterKey(e, () => {
                            newJenisBuku = addToArray(
                                detail.perpustakaan.jenis_buku,
                                newJenisBuku,
                            );
                        })}
                />
                <button
                    type="button"
                    class="btn btn-primary rounded-xl px-5 shadow-lg shadow-primary/20 transition-transform hover:-translate-y-0.5 min-h-11"
                    onclick={() => {
                        newJenisBuku = addToArray(
                            detail.perpustakaan.jenis_buku,
                            newJenisBuku,
                        );
                    }}
                >
                    Tambah
                </button>
            </div>
            <div class="flex flex-wrap gap-2">
                {#each detail.perpustakaan.jenis_buku as buku, index}
                    <div
                        class="badge font-medium px-3 py-3 rounded-xl bg-secondary/10 text-secondary border-secondary/20 gap-2"
                    >
                        {buku}
                        <button
                            title="Hapus Jenis Buku"
                            type="button"
                            class="btn btn-circle btn-ghost btn-xs h-5 w-5 min-h-0 p-0 hover:bg-secondary/20"
                            onclick={() =>
                                removeFromArray(
                                    detail.perpustakaan.jenis_buku,
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
        </div>
    </div>
</div>

<div
    class="card rounded-2xl border border-base-200 bg-base-100 shadow-xl shadow-base-200/40 overflow-hidden"
>
    <div class="card-body p-5 sm:p-6">
        <div class="flex items-center gap-3 mb-5">
            <div
                class="flex h-10 w-10 items-center justify-center rounded-xl bg-linear-to-br from-success/20 to-primary/20 text-success font-bold shadow-inner"
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
                        d="M9 20l-5.447-2.724A1 1 0 013 16.382V5.618a1 1 0 011.447-.894L9 7m0 13l6-3m-6 3V7m6 10l4.553 2.276A1 1 0 0021 18.382V7.618a1 1 0 00-.553-.894L15 4m0 13V4m0 0L9 7"
                    />
                </svg>
            </div>
            <h2
                class="card-title text-xl font-extrabold tracking-tight text-base-content"
            >
                Kalibrasi Kiblat
            </h2>
        </div>

        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div class="form-control">
                <label class="label py-2" for="kalibrasi_azimut">
                    <span
                        class="label-text font-bold text-base-content/80 text-xs uppercase tracking-wider"
                        >Azimut</span
                    >
                </label>
                <input
                    id="kalibrasi_azimut"
                    type="text"
                    class="input input-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-11 font-mono"
                    placeholder="295.12"
                    bind:value={detail.kalibrasi_arah_kiblat.azimut}
                />
            </div>
            <div class="form-control">
                <label class="label py-2" for="kalibrasi_tanggal">
                    <span
                        class="label-text font-bold text-base-content/80 text-xs uppercase tracking-wider"
                        >Tanggal Kalibrasi</span
                    >
                </label>
                <input
                    id="kalibrasi_tanggal"
                    type="date"
                    class="input input-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-11"
                    bind:value={detail.kalibrasi_arah_kiblat.tanggal_kalibrasi}
                />
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
                        d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z"
                    />
                </svg>
            </div>
            <h2
                class="card-title text-xl font-extrabold tracking-tight text-base-content"
            >
                Fasilitas & Kegiatan
            </h2>
        </div>

        <div class="mb-5">
            <label class="label py-2" for="fu_in">
                <span
                    class="label-text font-bold text-base-content/80 text-xs uppercase tracking-wider"
                    >Fasilitas Umum</span
                >
            </label>
            <div class="mb-3 flex gap-2">
                <input
                    id="fu_in"
                    type="text"
                    class="input input-bordered flex-1 rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-11"
                    placeholder="Contoh: Area Parkir"
                    bind:value={newFasilitasUmum}
                    onkeydown={(e) =>
                        handleEnterKey(e, () => {
                            newFasilitasUmum = addToArray(
                                detail.fasilitas_umum,
                                newFasilitasUmum,
                            );
                        })}
                />
                <button
                    type="button"
                    class="btn btn-primary rounded-xl px-5 shadow-lg shadow-primary/20 transition-transform hover:-translate-y-0.5 min-h-11"
                    onclick={() => {
                        newFasilitasUmum = addToArray(
                            detail.fasilitas_umum,
                            newFasilitasUmum,
                        );
                    }}
                >
                    Tambah
                </button>
            </div>
            <div class="flex flex-wrap gap-2">
                {#each detail.fasilitas_umum as f, i}
                    <div
                        class="badge font-medium px-3 py-3 rounded-xl bg-secondary/10 text-secondary border-secondary/20 gap-2"
                    >
                        {f}
                        <button
                            title="Hapus Fasilitas Umum"
                            type="button"
                            class="btn btn-circle btn-ghost btn-xs h-5 w-5 min-h-0 p-0 hover:bg-secondary/20"
                            onclick={() =>
                                removeFromArray(detail.fasilitas_umum, i)}
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
        </div>

        <div class="mb-5">
            <label class="label py-2" for="fra_in">
                <span
                    class="label-text font-bold text-base-content/80 text-xs uppercase tracking-wider"
                    >Fasilitas Ramah Anak</span
                >
            </label>
            <div class="mb-3 flex gap-2">
                <input
                    id="fra_in"
                    type="text"
                    class="input input-bordered flex-1 rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-11"
                    placeholder="Contoh: Ruang Bermain"
                    bind:value={newFasilitasRamanAnak}
                    onkeydown={(e) =>
                        handleEnterKey(e, () => {
                            newFasilitasRamanAnak = addToArray(
                                detail.fasilitas_raman_anak,
                                newFasilitasRamanAnak,
                            );
                        })}
                />
                <button
                    type="button"
                    class="btn btn-primary rounded-xl px-5 shadow-lg shadow-primary/20 transition-transform hover:-translate-y-0.5 min-h-11"
                    onclick={() => {
                        newFasilitasRamanAnak = addToArray(
                            detail.fasilitas_raman_anak,
                            newFasilitasRamanAnak,
                        );
                    }}
                >
                    Tambah
                </button>
            </div>
            <div class="flex flex-wrap gap-2">
                {#each detail.fasilitas_raman_anak as f, i}
                    <div
                        class="badge font-medium px-3 py-3 rounded-xl bg-accent/10 text-accent border-accent/20 gap-2"
                    >
                        {f}
                        <button
                            title="Hapus Fasilitas Raman Anak"
                            type="button"
                            class="btn btn-circle btn-ghost btn-xs h-5 w-5 min-h-0 p-0 hover:bg-accent/20"
                            onclick={() =>
                                removeFromArray(detail.fasilitas_raman_anak, i)}
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
        </div>

        <div class="mb-5">
            <label class="label py-2" for="fd_in">
                <span
                    class="label-text font-bold text-base-content/80 text-xs uppercase tracking-wider"
                    >Fasilitas Disabilitas</span
                >
            </label>
            <div class="mb-3 flex gap-2">
                <input
                    id="fd_in"
                    type="text"
                    class="input input-bordered flex-1 rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-11"
                    placeholder="Contoh: Jalur Kursi Roda"
                    bind:value={newFasilitasDisabilitas}
                    onkeydown={(e) =>
                        handleEnterKey(e, () => {
                            newFasilitasDisabilitas = addToArray(
                                detail.fasilitas_disabilitas,
                                newFasilitasDisabilitas,
                            );
                        })}
                />
                <button
                    type="button"
                    class="btn btn-primary rounded-xl px-5 shadow-lg shadow-primary/20 transition-transform hover:-translate-y-0.5 min-h-11"
                    onclick={() => {
                        newFasilitasDisabilitas = addToArray(
                            detail.fasilitas_disabilitas,
                            newFasilitasDisabilitas,
                        );
                    }}
                >
                    Tambah
                </button>
            </div>
            <div class="flex flex-wrap gap-2">
                {#each detail.fasilitas_disabilitas as f, i}
                    <div
                        class="badge font-medium px-3 py-3 rounded-xl bg-info/10 text-info border-info/20 gap-2"
                    >
                        {f}
                        <button
                            title="Hapus Fasilitas Disabilitas"
                            type="button"
                            class="btn btn-circle btn-ghost btn-xs h-5 w-5 min-h-0 p-0 hover:bg-info/20"
                            onclick={() =>
                                removeFromArray(
                                    detail.fasilitas_disabilitas,
                                    i,
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
        </div>

        <div>
            <label class="label py-2" for="keg_in">
                <span
                    class="label-text font-bold text-base-content/80 text-xs uppercase tracking-wider"
                    >Kegiatan Rutin</span
                >
            </label>
            <div class="mb-3 flex gap-2">
                <input
                    id="keg_in"
                    type="text"
                    class="input input-bordered flex-1 rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-11"
                    placeholder="Contoh: Pengajian Mingguan"
                    bind:value={newKegiatan}
                    onkeydown={(e) =>
                        handleEnterKey(e, () => {
                            newKegiatan = addToArray(
                                detail.kegiatan,
                                newKegiatan,
                            );
                        })}
                />
                <button
                    type="button"
                    class="btn btn-primary rounded-xl px-5 shadow-lg shadow-primary/20 transition-transform hover:-translate-y-0.5 min-h-11"
                    onclick={() => {
                        newKegiatan = addToArray(detail.kegiatan, newKegiatan);
                    }}
                >
                    Tambah
                </button>
            </div>
            <div class="flex flex-wrap gap-2">
                {#each detail.kegiatan as k, i}
                    <div
                        class="badge font-medium px-3 py-3 rounded-xl bg-success/10 text-success border-success/20 gap-2"
                    >
                        {k}
                        <button
                            title="Hapus Kegiatan"
                            type="button"
                            class="btn btn-circle btn-ghost btn-xs h-5 w-5 min-h-0 p-0 hover:bg-success/20"
                            onclick={() => removeFromArray(detail.kegiatan, i)}
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
        </div>
    </div>
</div>
