<script lang="ts">
    interface Props {
        detail?: Record<string, any>;
    }

    let { detail = {} }: Props = $props();
</script>

<div class="flex flex-col gap-6">
    {#if detail.kalibrasi_arah_kiblat}
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
                            ><path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M9 20l-5.447-2.724A1 1 0 013 16.382V5.618a1 1 0 011.447-.894L9 7m0 13l6-3m-6 3V7m6 10l4.553 2.276A1 1 0 0021 18.382V7.618a1 1 0 00-.553-.894L15 4m0 13V4m0 0L9 7"
                            /></svg
                        >
                    </div>
                    <h2
                        class="card-title text-xl font-extrabold tracking-tight text-base-content"
                    >
                        Kalibrasi Arah Kiblat
                    </h2>
                </div>
                <div class="grid grid-cols-1 gap-5 sm:grid-cols-2">
                    <div
                        class="rounded-xl bg-base-200/30 p-4 border border-base-200/50"
                    >
                        <p
                            class="mb-1 block text-xs font-bold text-base-content/80 uppercase tracking-wider"
                        >
                            Azimut
                        </p>
                        <p class="font-mono text-sm font-semibold text-primary">
                            {detail.kalibrasi_arah_kiblat.azimut || "-"}
                        </p>
                    </div>
                    <div
                        class="rounded-xl bg-base-200/30 p-4 border border-base-200/50"
                    >
                        <p
                            class="mb-1 block text-xs font-bold text-base-content/80 uppercase tracking-wider"
                        >
                            Tanggal Kalibrasi
                        </p>
                        <p class="font-medium text-sm text-base-content">
                            {detail.kalibrasi_arah_kiblat.tanggal_kalibrasi ||
                                "-"}
                        </p>
                    </div>
                </div>
            </div>
        </div>
    {/if}

    {#if detail.sdm_masjid}
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
                            ><path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
                            /></svg
                        >
                    </div>
                    <h2
                        class="card-title text-xl font-extrabold tracking-tight text-base-content"
                    >
                        SDM Mushola
                    </h2>
                </div>
                <div class="grid grid-cols-1 gap-5 sm:grid-cols-3">
                    <div
                        class="flex flex-col items-center justify-center rounded-2xl bg-base-200/30 p-4 border border-base-200/50 hover:bg-base-200/50 transition-colors"
                    >
                        <div
                            class="text-xs font-bold text-base-content/60 uppercase tracking-wider mb-2"
                        >
                            Imam
                        </div>
                        <div class="text-3xl font-extrabold text-primary">
                            {detail.sdm_masjid.jumlah_imam || 0}
                        </div>
                    </div>
                    <div
                        class="flex flex-col items-center justify-center rounded-2xl bg-base-200/30 p-4 border border-base-200/50 hover:bg-base-200/50 transition-colors"
                    >
                        <div
                            class="text-xs font-bold text-base-content/60 uppercase tracking-wider mb-2"
                        >
                            Muadzin
                        </div>
                        <div class="text-3xl font-extrabold text-secondary">
                            {detail.sdm_masjid.jumlah_muadzin || 0}
                        </div>
                    </div>
                    <div
                        class="flex flex-col items-center justify-center rounded-2xl bg-base-200/30 p-4 border border-base-200/50 hover:bg-base-200/50 transition-colors"
                    >
                        <div
                            class="text-xs font-bold text-base-content/60 uppercase tracking-wider mb-2"
                        >
                            Jemaah
                        </div>
                        <div class="text-3xl font-extrabold text-success">
                            {detail.sdm_masjid.jumlah_jemaah || 0}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    {/if}

    {#if detail.fasilitas_umum && detail.fasilitas_umum.length > 0}
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
                            ><path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z"
                            /></svg
                        >
                    </div>
                    <h2
                        class="card-title text-xl font-extrabold tracking-tight text-base-content"
                    >
                        Fasilitas Umum
                    </h2>
                </div>
                <div class="flex flex-wrap gap-2">
                    {#each detail.fasilitas_umum as fasilitas (fasilitas)}
                        <span
                            class="badge font-medium px-3 py-3 rounded-xl bg-secondary/10 text-secondary border-secondary/20"
                            >{fasilitas}</span
                        >
                    {/each}
                </div>
            </div>
        </div>
    {/if}

    {#if detail.nama_imam && detail.nama_imam.length > 0}
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
                            ><path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
                            /></svg
                        >
                    </div>
                    <h2
                        class="card-title text-xl font-extrabold tracking-tight text-base-content"
                    >
                        Nama Imam
                    </h2>
                </div>
                <div class="flex flex-wrap gap-2">
                    {#each detail.nama_imam as imam (imam)}
                        <span
                            class="badge font-medium px-3 py-3 rounded-xl bg-primary/10 text-primary border-primary/20"
                            >{imam}</span
                        >
                    {/each}
                </div>
            </div>
        </div>
    {/if}

    {#if detail.nama_muadzin && detail.nama_muadzin.length > 0}
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
                            ><path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M11 5.882V19.24a1.76 1.76 0 01-3.417.592l-2.147-6.15M18 13a3 3 0 100-6M5.436 13.683A4.001 4.001 0 017 6h1.832c4.1 0 7.625-1.234 9.168-3v14c-1.543-1.766-5.067-3-9.168-3H7a3.988 3.988 0 01-1.564-.317z"
                            /></svg
                        >
                    </div>
                    <h2
                        class="card-title text-xl font-extrabold tracking-tight text-base-content"
                    >
                        Nama Muadzin
                    </h2>
                </div>
                <div class="flex flex-wrap gap-2">
                    {#each detail.nama_muadzin as muadzin (muadzin)}
                        <span
                            class="badge font-medium px-3 py-3 rounded-xl bg-info/10 text-info border-info/20"
                            >{muadzin}</span
                        >
                    {/each}
                </div>
            </div>
        </div>
    {/if}

    {#if detail.kegiatan && detail.kegiatan.length > 0}
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
                            ><path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                            /></svg
                        >
                    </div>
                    <h2
                        class="card-title text-xl font-extrabold tracking-tight text-base-content"
                    >
                        Kegiatan
                    </h2>
                </div>
                <div class="flex flex-wrap gap-2">
                    {#each detail.kegiatan as kegiatan (kegiatan)}
                        <span
                            class="badge font-medium px-3 py-3 rounded-xl bg-success/10 text-success border-success/20"
                            >{kegiatan}</span
                        >
                    {/each}
                </div>
            </div>
        </div>
    {/if}
</div>
