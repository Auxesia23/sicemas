<script lang="ts">
    import { DEFAULT_DETAIL_MUSHOLA } from "$lib/schemas/siteInput";
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

    let newFasilitasUmum = $state("");
    let newNamaImam = $state("");
    let newNamaMuazdin = $state("");
    let newKegiatan = $state("");

    if (!detail?.kalibrasi_arah_kiblat) {
        detail = structuredClone(DEFAULT_DETAIL_MUSHOLA);
    }
</script>

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
                        >Tanggal</span
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
                SDM Mushola <span class="text-error">*</span>
            </h2>
        </div>

        {#if errors?.sdm_masjid && typeof errors.sdm_masjid === "string"}
            <div
                class="mb-4 text-sm font-medium text-error bg-error/10 rounded-xl p-3 border border-error/20"
            >
                {errors.sdm_masjid}
            </div>
        {/if}

        <div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
            {#each [{ id: "imam", label: "Imam", path: "jumlah_imam" }, { id: "muadzin", label: "Muadzin", path: "jumlah_muadzin" }, { id: "jemaah", label: "Jemaah", path: "jumlah_jemaah" }] as item}
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
                        placeholder="0"
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
                        d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z"
                    />
                </svg>
            </div>
            <h2
                class="card-title text-xl font-extrabold tracking-tight text-base-content"
            >
                Fasilitas Umum
            </h2>
        </div>

        <div>
            <label class="label py-2" for="fasilitas_umum_input">
                <span
                    class="label-text font-bold text-base-content/80 text-xs uppercase tracking-wider"
                    >Tambah Fasilitas</span
                >
            </label>
            <div class="mb-3 flex gap-2">
                <input
                    id="fasilitas_umum_input"
                    type="text"
                    class="input input-bordered flex-1 rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-11"
                    placeholder="Tambah fasilitas..."
                    bind:value={newFasilitasUmum}
                    onkeydown={(e) =>
                        handleEnterKey(
                            e,
                            () =>
                                (newFasilitasUmum = addToArray(
                                    detail.fasilitas_umum,
                                    newFasilitasUmum,
                                )),
                        )}
                />
                <button
                    type="button"
                    class="btn btn-primary rounded-xl px-5 shadow-lg shadow-primary/20 transition-transform hover:-translate-y-0.5 min-h-11"
                    onclick={() =>
                        (newFasilitasUmum = addToArray(
                            detail.fasilitas_umum,
                            newFasilitasUmum,
                        ))}
                >
                    Tambah
                </button>
            </div>
            <div class="flex flex-wrap gap-2">
                {#each detail.fasilitas_umum || [] as fasilitas, index (index)}
                    <div
                        class="badge font-medium px-3 py-3 rounded-xl bg-secondary/10 text-secondary border-secondary/20 gap-2"
                    >
                        {fasilitas}
                        <button
                            type="button"
                            class="btn btn-circle btn-ghost btn-xs h-5 w-5 min-h-0 p-0 hover:bg-secondary/20"
                            aria-label="Hapus fasilitas {fasilitas}"
                            onclick={() =>
                                removeFromArray(detail.fasilitas_umum, index)}
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
                {#if detail.fasilitas_umum?.length === 0}
                    <span class="text-sm font-medium text-base-content/50"
                        >Belum ada fasilitas umum</span
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
                Pengurus & Kegiatan
            </h2>
        </div>

        <div class="mb-5">
            <label class="label py-2" for="nama_imam_input">
                <span
                    class="label-text font-bold text-base-content/80 text-xs uppercase tracking-wider"
                >
                    Nama Imam <span class="text-error">*</span>
                </span>
            </label>
            <div class="mb-3 flex gap-2">
                <input
                    id="nama_imam_input"
                    type="text"
                    class="input input-bordered flex-1 rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-11 {errors?.nama_imam
                        ? 'input-error'
                        : ''}"
                    placeholder="Tambah nama imam..."
                    bind:value={newNamaImam}
                    onkeydown={(e) =>
                        handleEnterKey(e, () => {
                            newNamaImam = addToArray(
                                detail.nama_imam,
                                newNamaImam,
                            );
                            clearNestedError(errors, "nama_imam");
                        })}
                />
                <button
                    type="button"
                    class="btn btn-primary rounded-xl px-5 shadow-lg shadow-primary/20 transition-transform hover:-translate-y-0.5 min-h-11"
                    onclick={() => {
                        newNamaImam = addToArray(detail.nama_imam, newNamaImam);
                        clearNestedError(errors, "nama_imam");
                    }}
                >
                    Tambah
                </button>
            </div>
            <div class="flex flex-wrap gap-2">
                {#each detail.nama_imam || [] as imam, index (index)}
                    <div
                        class="badge font-medium px-3 py-3 rounded-xl bg-primary/10 text-primary border-primary/20 gap-2"
                    >
                        {imam}
                        <button
                            type="button"
                            class="btn btn-circle btn-ghost btn-xs h-5 w-5 min-h-0 p-0 hover:bg-primary/20"
                            aria-label="Hapus imam {imam}"
                            onclick={() =>
                                removeFromArray(detail.nama_imam, index)}
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
            {#if errors?.nama_imam}
                <span class="mt-1 block text-xs text-error font-medium"
                    >{errors.nama_imam}</span
                >
            {/if}
        </div>

        <div class="mb-5">
            <label class="label py-2" for="nama_muazdin_input">
                <span
                    class="label-text font-bold text-base-content/80 text-xs uppercase tracking-wider"
                >
                    Nama Muadzin <span class="text-error">*</span>
                </span>
            </label>
            <div class="mb-3 flex gap-2">
                <input
                    id="nama_muazdin_input"
                    type="text"
                    class="input input-bordered flex-1 rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-11 {errors?.nama_muadzin
                        ? 'input-error'
                        : ''}"
                    placeholder="Tambah nama muadzin..."
                    bind:value={newNamaMuazdin}
                    onkeydown={(e) =>
                        handleEnterKey(e, () => {
                            newNamaMuazdin = addToArray(
                                detail.nama_muadzin,
                                newNamaMuazdin,
                            );
                            clearNestedError(errors, "nama_muadzin");
                        })}
                />
                <button
                    type="button"
                    class="btn btn-primary rounded-xl px-5 shadow-lg shadow-primary/20 transition-transform hover:-translate-y-0.5 min-h-11"
                    onclick={() => {
                        newNamaMuazdin = addToArray(
                            detail.nama_muadzin,
                            newNamaMuazdin,
                        );
                        clearNestedError(errors, "nama_muadzin");
                    }}
                >
                    Tambah
                </button>
            </div>
            <div class="flex flex-wrap gap-2">
                {#each detail.nama_muadzin || [] as muadzin, index (index)}
                    <div
                        class="badge font-medium px-3 py-3 rounded-xl bg-info/10 text-info border-info/20 gap-2"
                    >
                        {muadzin}
                        <button
                            type="button"
                            class="btn btn-circle btn-ghost btn-xs h-5 w-5 min-h-0 p-0 hover:bg-info/20"
                            aria-label="Hapus muadzin {muadzin}"
                            onclick={() =>
                                removeFromArray(detail.nama_muadzin, index)}
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
            {#if errors?.nama_muadzin}
                <span class="mt-1 block text-xs text-error font-medium"
                    >{errors.nama_muadzin}</span
                >
            {/if}
        </div>

        <div>
            <label class="label py-2" for="kegiatan_input">
                <span
                    class="label-text font-bold text-base-content/80 text-xs uppercase tracking-wider"
                >
                    Kegiatan <span
                        class="text-xs font-medium text-base-content/50 ml-1"
                        >(Opsional)</span
                    >
                </span>
            </label>
            <div class="mb-3 flex gap-2">
                <input
                    id="kegiatan_input"
                    type="text"
                    class="input input-bordered flex-1 rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-11"
                    placeholder="Tambah kegiatan..."
                    bind:value={newKegiatan}
                    onkeydown={(e) =>
                        handleEnterKey(
                            e,
                            () =>
                                (newKegiatan = addToArray(
                                    detail.kegiatan,
                                    newKegiatan,
                                )),
                        )}
                />
                <button
                    type="button"
                    class="btn btn-primary rounded-xl px-5 shadow-lg shadow-primary/20 transition-transform hover:-translate-y-0.5 min-h-11"
                    onclick={() =>
                        (newKegiatan = addToArray(
                            detail.kegiatan,
                            newKegiatan,
                        ))}
                >
                    Tambah
                </button>
            </div>
            <div class="flex flex-wrap gap-2">
                {#each detail.kegiatan || [] as kegiatan, index (index)}
                    <div
                        class="badge font-medium px-3 py-3 rounded-xl bg-success/10 text-success border-success/20 gap-2"
                    >
                        {kegiatan}
                        <button
                            type="button"
                            class="btn btn-circle btn-ghost btn-xs h-5 w-5 min-h-0 p-0 hover:bg-success/20"
                            aria-label="Hapus kegiatan {kegiatan}"
                            onclick={() =>
                                removeFromArray(detail.kegiatan, index)}
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
                {#if detail.kegiatan?.length === 0}
                    <span class="text-sm font-medium text-base-content/50"
                        >Belum ada kegiatan</span
                    >
                {/if}
            </div>
        </div>
    </div>
</div>
