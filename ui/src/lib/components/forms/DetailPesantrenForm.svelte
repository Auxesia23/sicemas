<script lang="ts">
    import { DEFAULT_DETAIL_PESANTREN } from "$lib/schemas/siteInput";
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

    let newFasilitasPondok = $state("");
    let newSekretaris = $state("");
    let newBendahara = $state("");

    if (!detail?.nama_yayasan && detail?.nama_yayasan !== "") {
        detail = structuredClone(DEFAULT_DETAIL_PESANTREN);
    }

    function addSekretaris() {
        newSekretaris = addToArray(
            detail.kepengurusan.sekretaris,
            newSekretaris,
        );
        clearNestedError(errors, "kepengurusan", "sekretaris");
    }

    function addBendahara() {
        newBendahara = addToArray(detail.kepengurusan.bendahara, newBendahara);
        clearNestedError(errors, "kepengurusan", "bendahara");
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
                        d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
                    />
                </svg>
            </div>
            <h2
                class="card-title text-xl font-extrabold tracking-tight text-base-content"
            >
                Informasi Yayasan <span class="text-error">*</span>
            </h2>
        </div>

        <div class="flex flex-col gap-4">
            <div class="form-control">
                <label class="label py-2" for="nama_yayasan">
                    <span
                        class="label-text font-bold text-base-content/80 text-xs uppercase tracking-wider"
                    >
                        Nama Yayasan <span
                            class="text-xs font-medium text-base-content/50 ml-1"
                            >(Opsional)</span
                        >
                    </span>
                </label>
                <input
                    id="nama_yayasan"
                    type="text"
                    class="input input-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-11"
                    bind:value={detail.nama_yayasan}
                    oninput={() => clearNestedError(errors, "nama_yayasan")}
                />
            </div>
            <div class="form-control">
                <label class="label py-2" for="pimpinan_pesantren">
                    <span
                        class="label-text font-bold text-base-content/80 text-xs uppercase tracking-wider"
                    >
                        Pimpinan Pesantren <span class="text-error">*</span>
                    </span>
                </label>
                <input
                    id="pimpinan_pesantren"
                    type="text"
                    class="input input-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-11 {errors?.pimpinan_pesantren
                        ? 'input-error'
                        : ''}"
                    bind:value={detail.pimpinan_pesantren}
                    oninput={() =>
                        clearNestedError(errors, "pimpinan_pesantren")}
                />
                {#if errors?.pimpinan_pesantren}
                    <span class="mt-1 text-xs text-error font-medium"
                        >{errors.pimpinan_pesantren}</span
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
                        d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
                    />
                </svg>
            </div>
            <h2
                class="card-title text-xl font-extrabold tracking-tight text-base-content"
            >
                Kepengurusan <span class="text-error">*</span>
            </h2>
        </div>

        <div class="form-control mb-5">
            <label class="label py-2" for="kepengurusan_ketua">
                <span
                    class="label-text font-bold text-base-content/80 text-xs uppercase tracking-wider"
                >
                    Ketua Yayasan/Pesantren <span class="text-error">*</span>
                </span>
            </label>
            <input
                id="kepengurusan_ketua"
                type="text"
                class="input input-bordered w-full rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-11 {errors
                    ?.kepengurusan?.ketua
                    ? 'input-error'
                    : ''}"
                bind:value={detail.kepengurusan.ketua}
                oninput={() =>
                    clearNestedError(errors, "kepengurusan", "ketua")}
            />
            {#if errors?.kepengurusan?.ketua}
                <span class="mt-1 text-xs text-error font-medium"
                    >{errors.kepengurusan.ketua}</span
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
                            ?.kepengurusan?.sekretaris
                            ? 'input-error'
                            : ''}"
                        placeholder="Tambah Sekretaris..."
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
                    {#each detail.kepengurusan.sekretaris as p, i}
                        <div
                            class="badge font-medium px-3 py-3 rounded-xl bg-info/10 text-info border-info/20 gap-2"
                        >
                            {p}
                            <button
                                type="button"
                                class="btn btn-circle btn-ghost btn-xs h-5 w-5 min-h-0 p-0 hover:bg-info/20"
                                aria-label="Hapus sekretaris {p}"
                                onclick={() =>
                                    removeFromArray(
                                        detail.kepengurusan.sekretaris,
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
                {#if errors?.kepengurusan?.sekretaris}
                    <span class="mt-1 text-xs text-error font-medium"
                        >{errors.kepengurusan.sekretaris}</span
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
                            ?.kepengurusan?.bendahara
                            ? 'input-error'
                            : ''}"
                        placeholder="Tambah Bendahara..."
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
                    {#each detail.kepengurusan.bendahara as p, i}
                        <div
                            class="badge font-medium px-3 py-3 rounded-xl bg-warning/10 text-warning border-warning/20 gap-2"
                        >
                            {p}
                            <button
                                type="button"
                                class="btn btn-circle btn-ghost btn-xs h-5 w-5 min-h-0 p-0 hover:bg-warning/20"
                                aria-label="Hapus bendahara {p}"
                                onclick={() =>
                                    removeFromArray(
                                        detail.kepengurusan.bendahara,
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
                {#if errors?.kepengurusan?.bendahara}
                    <span class="mt-1 text-xs text-error font-medium"
                        >{errors.kepengurusan.bendahara}</span
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
                        d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"
                    />
                </svg>
            </div>
            <h2
                class="card-title text-xl font-extrabold tracking-tight text-base-content"
            >
                Data Santri <span class="text-error">*</span>
            </h2>
        </div>

        <div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
            {#each ["mondok", "tidak_mondok", "total"] as cat}
                <div
                    class="card rounded-xl bg-base-200/30 border border-base-200/50 overflow-hidden"
                >
                    <div class="card-body p-4">
                        <h3
                            class="card-title text-sm font-bold uppercase tracking-wider text-base-content/80 mb-3"
                        >
                            {cat.replace("_", " ")}
                        </h3>
                        <div class="grid grid-cols-2 gap-3">
                            <div class="form-control">
                                <label class="label py-1" for="{cat}_pria">
                                    <span
                                        class="label-text text-xs font-bold text-base-content/60 uppercase tracking-wider"
                                        >Pria</span
                                    >
                                </label>
                                <input
                                    id="{cat}_pria"
                                    type="number"
                                    min="0"
                                    class="input input-bordered input-sm w-full rounded-lg bg-base-100 {errors
                                        ?.santri?.[cat]?.pria
                                        ? 'input-error'
                                        : ''}"
                                    bind:value={detail.santri[cat].pria}
                                    oninput={() =>
                                        clearNestedError(
                                            errors,
                                            "santri",
                                            cat,
                                            "pria",
                                        )}
                                />
                                {#if errors?.santri?.[cat]?.pria}
                                    <span
                                        class="mt-0.5 text-[10px] text-error font-medium"
                                        >{errors.santri[cat].pria}</span
                                    >
                                {/if}
                            </div>
                            <div class="form-control">
                                <label class="label py-1" for="{cat}_wanita">
                                    <span
                                        class="label-text text-xs font-bold text-base-content/60 uppercase tracking-wider"
                                        >Wanita</span
                                    >
                                </label>
                                <input
                                    id="{cat}_wanita"
                                    type="number"
                                    min="0"
                                    class="input input-bordered input-sm w-full rounded-lg bg-base-100 {errors
                                        ?.santri?.[cat]?.wanita
                                        ? 'input-error'
                                        : ''}"
                                    bind:value={detail.santri[cat].wanita}
                                    oninput={() =>
                                        clearNestedError(
                                            errors,
                                            "santri",
                                            cat,
                                            "wanita",
                                        )}
                                />
                                {#if errors?.santri?.[cat]?.wanita}
                                    <span
                                        class="mt-0.5 text-[10px] text-error font-medium"
                                        >{errors.santri[cat].wanita}</span
                                    >
                                {/if}
                            </div>
                        </div>
                    </div>
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
                        d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z"
                    />
                </svg>
            </div>
            <h2
                class="card-title text-xl font-extrabold tracking-tight text-base-content"
            >
                Fasilitas Pondok
            </h2>
        </div>

        <div>
            <label class="label py-2" for="fasilitas_pondok_input">
                <span
                    class="label-text font-bold text-base-content/80 text-xs uppercase tracking-wider"
                    >Tambah Fasilitas</span
                >
            </label>
            <div class="mb-3 flex gap-2">
                <input
                    id="fasilitas_pondok_input"
                    type="text"
                    class="input input-bordered flex-1 rounded-xl bg-base-200/30 focus:bg-base-100 transition-colors min-h-11"
                    placeholder="Tambah fasilitas..."
                    bind:value={newFasilitasPondok}
                    onkeydown={(e) =>
                        handleEnterKey(
                            e,
                            () =>
                                (newFasilitasPondok = addToArray(
                                    detail.fasilitas_pondok,
                                    newFasilitasPondok,
                                )),
                        )}
                />
                <button
                    type="button"
                    class="btn btn-primary rounded-xl px-5 shadow-lg shadow-primary/20 transition-transform hover:-translate-y-0.5 min-h-11"
                    onclick={() =>
                        (newFasilitasPondok = addToArray(
                            detail.fasilitas_pondok,
                            newFasilitasPondok,
                        ))}
                >
                    Tambah
                </button>
            </div>
            <div class="flex flex-wrap gap-2">
                {#each detail.fasilitas_pondok || [] as fasilitas, index (index)}
                    <div
                        class="badge font-medium px-3 py-3 rounded-xl bg-primary/10 text-primary border-primary/20 gap-2"
                    >
                        {fasilitas}
                        <button
                            type="button"
                            class="btn btn-circle btn-ghost btn-xs h-5 w-5 min-h-0 p-0 hover:bg-primary/20"
                            aria-label="Hapus fasilitas {fasilitas}"
                            onclick={() =>
                                removeFromArray(detail.fasilitas_pondok, index)}
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
                {#if detail.fasilitas_pondok?.length === 0}
                    <span class="text-sm font-medium text-base-content/50"
                        >Belum ada fasilitas pondok</span
                    >
                {/if}
            </div>
        </div>
    </div>
</div>
