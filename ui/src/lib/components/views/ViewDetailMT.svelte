<script lang="ts">
    interface Props {
        detail?: Record<string, any>;
    }

    let { detail = {} }: Props = $props();
</script>

<div class="flex flex-col gap-6">
    {#if detail.pengurus}
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
                        Pengurus
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
                            {detail.pengurus.ketua || "-"}
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
                        {#if detail.pengurus.sekretaris && detail.pengurus.sekretaris.length > 0}
                            <div class="flex flex-wrap gap-2">
                                {#each detail.pengurus.sekretaris as sekretaris (sekretaris)}
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
                        {#if detail.pengurus.bendahara && detail.pengurus.bendahara.length > 0}
                            <div class="flex flex-wrap gap-2">
                                {#each detail.pengurus.bendahara as bendahara (bendahara)}
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
            </div>
        </div>
    {/if}

    {#if detail.kehadiran_pria}
        <div
            class="card rounded-2xl border border-base-200 bg-base-100 shadow-xl shadow-base-200/40 overflow-hidden"
        >
            <div class="card-body p-5 sm:p-6">
                <div class="flex items-center gap-3 mb-5">
                    <div
                        class="flex h-10 w-10 items-center justify-center rounded-xl bg-linear-to-br from-info/20 to-primary/20 text-primary font-bold shadow-inner"
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
                        Kehadiran Pria
                    </h2>
                </div>

                <div class="grid grid-cols-1 gap-5 sm:grid-cols-3">
                    <div
                        class="flex flex-col items-center justify-center rounded-2xl bg-base-200/30 p-4 border border-base-200/50 hover:bg-base-200/50 transition-colors"
                    >
                        <div
                            class="text-xs font-bold text-base-content/60 uppercase tracking-wider mb-2"
                        >
                            Dewasa
                        </div>
                        <div class="text-3xl font-extrabold text-primary">
                            {detail.kehadiran_pria.jumlah_dewasa ?? 0}
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
                        <div class="text-3xl font-extrabold text-info">
                            {detail.kehadiran_pria.jumlah_remaja ?? 0}
                        </div>
                    </div>
                    <div
                        class="flex flex-col justify-center rounded-2xl bg-base-200/30 p-4 border border-base-200/50"
                    >
                        <p
                            class="mb-2 block text-xs font-bold text-base-content/80 uppercase tracking-wider text-center sm:text-left"
                        >
                            Waktu Pengajian
                        </p>
                        <div class="flex justify-center sm:justify-start">
                            {#if detail.kehadiran_pria.waktu_pengajian}
                                <span
                                    class="badge font-medium px-4 py-4 rounded-xl bg-primary/10 text-primary border-primary/20 text-sm"
                                >
                                    {detail.kehadiran_pria.waktu_pengajian}
                                </span>
                            {:else}
                                <span class="font-medium text-base-content/50"
                                    >-</span
                                >
                            {/if}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    {/if}

    {#if detail.kehadiran_wanita}
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
                                d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"
                            /></svg
                        >
                    </div>
                    <h2
                        class="card-title text-xl font-extrabold tracking-tight text-base-content"
                    >
                        Kehadiran Wanita
                    </h2>
                </div>

                <div class="grid grid-cols-1 gap-5 sm:grid-cols-3">
                    <div
                        class="flex flex-col items-center justify-center rounded-2xl bg-base-200/30 p-4 border border-base-200/50 hover:bg-base-200/50 transition-colors"
                    >
                        <div
                            class="text-xs font-bold text-base-content/60 uppercase tracking-wider mb-2"
                        >
                            Dewasa
                        </div>
                        <div class="text-3xl font-extrabold text-secondary">
                            {detail.kehadiran_wanita.jumlah_dewasa ?? 0}
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
                        <div class="text-3xl font-extrabold text-accent">
                            {detail.kehadiran_wanita.jumlah_remaja ?? 0}
                        </div>
                    </div>
                    <div
                        class="flex flex-col justify-center rounded-2xl bg-base-200/30 p-4 border border-base-200/50"
                    >
                        <p
                            class="mb-2 block text-xs font-bold text-base-content/80 uppercase tracking-wider text-center sm:text-left"
                        >
                            Waktu Pengajian
                        </p>
                        <div class="flex justify-center sm:justify-start">
                            {#if detail.kehadiran_wanita.waktu_pengajian}
                                <span
                                    class="badge font-medium px-4 py-4 rounded-xl bg-secondary/10 text-secondary border-secondary/20 text-sm"
                                >
                                    {detail.kehadiran_wanita.waktu_pengajian}
                                </span>
                            {:else}
                                <span class="font-medium text-base-content/50"
                                    >-</span
                                >
                            {/if}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    {/if}

    {#if detail.penceramah && detail.penceramah.length > 0}
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
                            ><path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z"
                            /></svg
                        >
                    </div>
                    <h2
                        class="card-title text-xl font-extrabold tracking-tight text-base-content"
                    >
                        Daftar Penceramah
                    </h2>
                </div>

                <div class="space-y-3">
                    {#each detail.penceramah as penceramah (penceramah.nama)}
                        <div
                            class="flex items-center justify-between rounded-xl bg-base-200/30 border border-base-200/50 p-4 hover:bg-base-200/60 transition-colors"
                        >
                            <div class="flex-1">
                                <div
                                    class="text-base font-bold text-base-content mb-0.5"
                                >
                                    {penceramah.nama}
                                </div>
                                <div
                                    class="text-sm font-medium text-base-content/60 flex items-center gap-1.5"
                                >
                                    <svg
                                        xmlns="http://www.w3.org/2000/svg"
                                        class="h-4 w-4 opacity-50"
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
                                    {penceramah.materi ||
                                        "Materi belum ditentukan"}
                                </div>
                            </div>
                        </div>
                    {/each}
                </div>
            </div>
        </div>
    {/if}
</div>
