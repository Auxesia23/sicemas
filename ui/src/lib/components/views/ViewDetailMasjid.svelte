<script lang="ts">
    interface Props {
        detail?: Record<string, any>;
    }

    let { detail = {} }: Props = $props();
</script>

<div class="flex flex-col gap-6">
    {#if detail.perpustakaan}
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
                                d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"
                            /></svg
                        >
                    </div>
                    <h2
                        class="card-title text-xl font-extrabold tracking-tight text-base-content"
                    >
                        Perpustakaan
                    </h2>
                </div>

                <div class="grid grid-cols-1 gap-5 sm:grid-cols-3">
                    <div
                        class="rounded-xl bg-base-200/30 p-3 border border-base-200/50"
                    >
                        <p
                            class="mb-1 block text-xs font-bold text-base-content/80 uppercase tracking-wider"
                        >
                            Luas
                        </p>
                        <p class="font-medium text-sm text-base-content">
                            {detail.perpustakaan.luas_m2
                                ? `${detail.perpustakaan.luas_m2} m²`
                                : "-"}
                        </p>
                    </div>
                    <div
                        class="rounded-xl bg-base-200/30 p-3 border border-base-200/50"
                    >
                        <p
                            class="mb-1 block text-xs font-bold text-base-content/80 uppercase tracking-wider"
                        >
                            Jumlah Pengurus
                        </p>
                        <p class="font-medium text-sm text-base-content">
                            {detail.perpustakaan.jumlah_pengurus || "-"}
                        </p>
                    </div>
                    <div
                        class="rounded-xl bg-base-200/30 p-3 border border-base-200/50"
                    >
                        <p
                            class="mb-1 block text-xs font-bold text-base-content/80 uppercase tracking-wider"
                        >
                            Kondisi
                        </p>
                        <div class="mt-0.5">
                            <span
                                class={`badge badge-sm font-medium px-2.5 py-3 rounded-xl ${detail.perpustakaan.kondisi === "Baik" ? "bg-success/10 text-success border-success/20" : "bg-base-200/50 text-base-content/70 border-base-200"}`}
                            >
                                <div
                                    class={`w-1.5 h-1.5 rounded-full mr-1.5 ${detail.perpustakaan.kondisi === "Baik" ? "bg-success" : "bg-base-content/40"}`}
                                ></div>
                                {detail.perpustakaan.kondisi || "-"}
                            </span>
                        </div>
                    </div>
                </div>

                {#if detail.perpustakaan.jenis_buku && detail.perpustakaan.jenis_buku.length > 0}
                    <div class="mt-5">
                        <p
                            class="mb-2 block text-xs font-bold text-base-content/80 uppercase tracking-wider"
                        >
                            Jenis Buku
                        </p>
                        <div class="flex flex-wrap gap-2">
                            {#each detail.perpustakaan.jenis_buku as buku (buku)}
                                <span
                                    class="badge font-medium px-3 py-3 rounded-xl bg-primary/10 text-primary border-primary/20"
                                    >{buku}</span
                                >
                            {/each}
                        </div>
                    </div>
                {/if}
            </div>
        </div>
    {/if}

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
                        SDM Masjid
                    </h2>
                </div>
                <div
                    class="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-6"
                >
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
                            Khotib
                        </div>
                        <div class="text-3xl font-extrabold text-secondary">
                            {detail.sdm_masjid.jumlah_khotib || 0}
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
                        <div class="text-3xl font-extrabold text-accent">
                            {detail.sdm_masjid.jumlah_muadzin || 0}
                        </div>
                    </div>
                    <div
                        class="flex flex-col items-center justify-center rounded-2xl bg-base-200/30 p-4 border border-base-200/50 hover:bg-base-200/50 transition-colors"
                    >
                        <div
                            class="text-xs font-bold text-base-content/60 uppercase tracking-wider mb-2"
                        >
                            Pengurus
                        </div>
                        <div class="text-3xl font-extrabold text-info">
                            {detail.sdm_masjid.jumlah_pengurus || 0}
                        </div>
                    </div>
                    <div
                        class="flex flex-col items-center justify-center rounded-2xl bg-base-200/30 p-4 border border-base-200/50 hover:bg-base-200/50 transition-colors"
                    >
                        <div
                            class="text-xs font-bold text-base-content/60 uppercase tracking-wider mb-2"
                        >
                            Remaja
                        </div>
                        <div class="text-3xl font-extrabold text-warning">
                            {detail.sdm_masjid.jumlah_remaja || 0}
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

    {#if detail.pengurus_dkm}
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
                                d="M10 6H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V8a2 2 0 00-2-2h-5m-4 0V5a2 2 0 114 0v1m-4 0a2 2 0 104 0m-5 8a2 2 0 100-4 2 2 0 000 4zm0 0c1.306 0 2.417.835 2.83 2M9 14a3.001 3.001 0 00-2.83 2M15 11h3m-3 4h2"
                            /></svg
                        >
                    </div>
                    <h2
                        class="card-title text-xl font-extrabold tracking-tight text-base-content"
                    >
                        Pengurus DKM
                    </h2>
                </div>

                <div class="grid grid-cols-1 gap-5 sm:grid-cols-3">
                    <div
                        class="rounded-xl bg-base-200/30 p-4 border border-base-200/50"
                    >
                        <p
                            class="mb-1 block text-xs font-bold text-base-content/80 uppercase tracking-wider"
                        >
                            Ketua
                        </p>
                        <p class="font-medium text-sm text-base-content">
                            {detail.pengurus_dkm.ketua || "-"}
                        </p>
                    </div>
                    <div
                        class="rounded-xl bg-base-200/30 p-4 border border-base-200/50"
                    >
                        <p
                            class="mb-2 block text-xs font-bold text-base-content/80 uppercase tracking-wider"
                        >
                            Sekretaris
                        </p>
                        {#if detail.pengurus_dkm.sekretaris && detail.pengurus_dkm.sekretaris.length > 0}
                            <div class="flex flex-wrap gap-2">
                                {#each detail.pengurus_dkm.sekretaris as sekretaris (sekretaris)}
                                    <span
                                        class="badge font-medium px-3 py-3 rounded-xl bg-base-200/50 text-base-content/80 border-base-300"
                                        >{sekretaris}</span
                                    >
                                {/each}
                            </div>
                        {:else}
                            <p class="font-medium text-sm text-base-content/50">
                                -
                            </p>
                        {/if}
                    </div>
                    <div
                        class="rounded-xl bg-base-200/30 p-4 border border-base-200/50"
                    >
                        <p
                            class="mb-2 block text-xs font-bold text-base-content/80 uppercase tracking-wider"
                        >
                            Bendahara
                        </p>
                        {#if detail.pengurus_dkm.bendahara && detail.pengurus_dkm.bendahara.length > 0}
                            <div class="flex flex-wrap gap-2">
                                {#each detail.pengurus_dkm.bendahara as bendahara (bendahara)}
                                    <span
                                        class="badge font-medium px-3 py-3 rounded-xl bg-base-200/50 text-base-content/80 border-base-300"
                                        >{bendahara}</span
                                    >
                                {/each}
                            </div>
                        {:else}
                            <p class="font-medium text-sm text-base-content/50">
                                -
                            </p>
                        {/if}
                    </div>
                </div>

                <div class="grid grid-cols-1 gap-5 sm:grid-cols-2 mt-5">
                    {#if detail.pengurus_dkm.nama_imam && detail.pengurus_dkm.nama_imam.length > 0}
                        <div>
                            <p
                                class="mb-2 block text-xs font-bold text-base-content/80 uppercase tracking-wider"
                            >
                                Nama Imam
                            </p>
                            <div class="flex flex-wrap gap-2">
                                {#each detail.pengurus_dkm.nama_imam as imam (imam)}
                                    <span
                                        class="badge font-medium px-3 py-3 rounded-xl bg-primary/10 text-primary border-primary/20"
                                        >{imam}</span
                                    >
                                {/each}
                            </div>
                        </div>
                    {/if}

                    {#if detail.pengurus_dkm.nama_muazdin && detail.pengurus_dkm.nama_muazdin.length > 0}
                        <div>
                            <p
                                class="mb-2 block text-xs font-bold text-base-content/80 uppercase tracking-wider"
                            >
                                Nama Muadzin
                            </p>
                            <div class="flex flex-wrap gap-2">
                                {#each detail.pengurus_dkm.nama_muazdin as muazdin (muazdin)}
                                    <span
                                        class="badge font-medium px-3 py-3 rounded-xl bg-secondary/10 text-secondary border-secondary/20"
                                        >{muazdin}</span
                                    >
                                {/each}
                            </div>
                        </div>
                    {/if}
                </div>
            </div>
        </div>
    {/if}

    {#if detail.fasilitas_umum || detail.fasilitas_raman_anak || detail.fasilitas_disabilitas}
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
                                d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z"
                            /></svg
                        >
                    </div>
                    <h2
                        class="card-title text-xl font-extrabold tracking-tight text-base-content"
                    >
                        Fasilitas
                    </h2>
                </div>

                <div class="space-y-6">
                    <div>
                        <p
                            class="mb-2 block text-xs font-bold text-base-content/80 uppercase tracking-wider"
                        >
                            Fasilitas Umum
                        </p>
                        {#if detail.fasilitas_umum && detail.fasilitas_umum.length > 0}
                            <div class="flex flex-wrap gap-2">
                                {#each detail.fasilitas_umum as fasilitas (fasilitas)}
                                    <span
                                        class="badge font-medium px-3 py-3 rounded-xl bg-secondary/10 text-secondary border-secondary/20"
                                        >{fasilitas}</span
                                    >
                                {/each}
                            </div>
                        {:else}
                            <p
                                class="text-sm font-medium text-base-content/50 bg-base-200/30 p-3 rounded-xl border border-base-200/50 inline-block"
                            >
                                Tidak ada fasilitas umum.
                            </p>
                        {/if}
                    </div>

                    <div>
                        <p
                            class="mb-2 block text-xs font-bold text-base-content/80 uppercase tracking-wider"
                        >
                            Fasilitas Ramah Anak
                        </p>
                        {#if detail.fasilitas_raman_anak && detail.fasilitas_raman_anak.length > 0}
                            <div class="flex flex-wrap gap-2">
                                {#each detail.fasilitas_raman_anak as fasilitas (fasilitas)}
                                    <span
                                        class="badge font-medium px-3 py-3 rounded-xl bg-accent/10 text-accent border-accent/20"
                                        >{fasilitas}</span
                                    >
                                {/each}
                            </div>
                        {:else}
                            <p
                                class="text-sm font-medium text-base-content/50 bg-base-200/30 p-3 rounded-xl border border-base-200/50 inline-block"
                            >
                                Tidak ada fasilitas ramah anak.
                            </p>
                        {/if}
                    </div>

                    <div>
                        <p
                            class="mb-2 block text-xs font-bold text-base-content/80 uppercase tracking-wider"
                        >
                            Fasilitas Disabilitas
                        </p>
                        {#if detail.fasilitas_disabilitas && detail.fasilitas_disabilitas.length > 0}
                            <div class="flex flex-wrap gap-2">
                                {#each detail.fasilitas_disabilitas as fasilitas (fasilitas)}
                                    <span
                                        class="badge font-medium px-3 py-3 rounded-xl bg-info/10 text-info border-info/20"
                                        >{fasilitas}</span
                                    >
                                {/each}
                            </div>
                        {:else}
                            <p
                                class="text-sm font-medium text-base-content/50 bg-base-200/30 p-3 rounded-xl border border-base-200/50 inline-block"
                            >
                                Tidak ada fasilitas disabilitas.
                            </p>
                        {/if}
                    </div>
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
