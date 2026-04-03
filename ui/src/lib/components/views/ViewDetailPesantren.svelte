<script lang="ts">
    interface Props {
        detail?: Record<string, any>;
    }

    let { detail = {} }: Props = $props();
</script>

<div class="flex flex-col gap-6">
    {#if detail.nama_yayasan || detail.pimpinan_pesantren}
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
                                d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
                            /></svg
                        >
                    </div>
                    <h2
                        class="card-title text-xl font-extrabold tracking-tight text-base-content"
                    >
                        Informasi Yayasan
                    </h2>
                </div>

                <div class="grid grid-cols-1 gap-5 sm:grid-cols-2">
                    <div
                        class="rounded-xl bg-base-200/30 p-4 border border-base-200/50"
                    >
                        <p
                            class="mb-1 block text-xs font-bold text-base-content/80 uppercase tracking-wider"
                        >
                            Nama Yayasan
                        </p>
                        <p class="font-medium text-sm text-base-content">
                            {detail.nama_yayasan || "-"}
                        </p>
                    </div>
                    <div
                        class="rounded-xl bg-base-200/30 p-4 border border-base-200/50"
                    >
                        <p
                            class="mb-1 block text-xs font-bold text-base-content/80 uppercase tracking-wider"
                        >
                            Pimpinan Pesantren
                        </p>
                        <p class="font-medium text-sm text-base-content">
                            {detail.pimpinan_pesantren || "-"}
                        </p>
                    </div>
                </div>
            </div>
        </div>
    {/if}

    {#if detail.kepengurusan}
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
                        Kepengurusan
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
                            {detail.kepengurusan.ketua || "-"}
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
                        {#if detail.kepengurusan.sekretaris && detail.kepengurusan.sekretaris.length > 0}
                            <div class="flex flex-wrap gap-2">
                                {#each detail.kepengurusan.sekretaris as sekretaris (sekretaris)}
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
                        {#if detail.kepengurusan.bendahara && detail.kepengurusan.bendahara.length > 0}
                            <div class="flex flex-wrap gap-2">
                                {#each detail.kepengurusan.bendahara as bendahara (bendahara)}
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

    {#if detail.santri}
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
                                d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"
                            /></svg
                        >
                    </div>
                    <h2
                        class="card-title text-xl font-extrabold tracking-tight text-base-content"
                    >
                        Data Santri
                    </h2>
                </div>

                <div
                    class="overflow-x-auto rounded-xl border border-base-200/50"
                >
                    <table class="table w-full">
                        <thead>
                            <tr
                                class="bg-base-200/30 text-xs uppercase tracking-wider font-bold text-base-content/70 border-b-2 border-base-200"
                            >
                                <th class="py-3 px-4">Kategori</th>
                                <th class="py-3 px-4 text-center">Pria</th>
                                <th class="py-3 px-4 text-center">Wanita</th>
                            </tr>
                        </thead>
                        <tbody class="text-sm">
                            <tr
                                class="border-b border-base-200/50 hover:bg-base-200/20 transition-colors"
                            >
                                <td
                                    class="py-3 px-4 font-medium text-base-content"
                                    >Mondok</td
                                >
                                <td
                                    class="py-3 px-4 text-center font-semibold text-primary"
                                >
                                    {detail.santri.mondok?.pria ?? 0}
                                </td>
                                <td
                                    class="py-3 px-4 text-center font-semibold text-secondary"
                                >
                                    {detail.santri.mondok?.wanita ?? 0}
                                </td>
                            </tr>
                            <tr
                                class="border-b border-base-200/50 hover:bg-base-200/20 transition-colors"
                            >
                                <td
                                    class="py-3 px-4 font-medium text-base-content"
                                    >Tidak Mondok</td
                                >
                                <td
                                    class="py-3 px-4 text-center font-semibold text-primary"
                                >
                                    {detail.santri.tidak_mondok?.pria ?? 0}
                                </td>
                                <td
                                    class="py-3 px-4 text-center font-semibold text-secondary"
                                >
                                    {detail.santri.tidak_mondok?.wanita ?? 0}
                                </td>
                            </tr>
                            <tr class="bg-base-200/20 font-bold">
                                <td class="py-3 px-4 text-base-content"
                                    >Total</td
                                >
                                <td class="py-3 px-4 text-center text-primary">
                                    {detail.santri.total?.pria ?? 0}
                                </td>
                                <td
                                    class="py-3 px-4 text-center text-secondary"
                                >
                                    {detail.santri.total?.wanita ?? 0}
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    {/if}

    {#if detail.fasilitas_pondok && detail.fasilitas_pondok.length > 0}
        <div
            class="card rounded-2xl border border-base-200 bg-base-100 shadow-xl shadow-base-200/40 overflow-hidden"
        >
            <div class="card-body p-5 sm:p-6">
                <div class="flex items-center gap-3 mb-5">
                    <div
                        class="flex h-10 w-10 items-center justify-center rounded-xl bg-linear-to-br from-success/20 to-info/20 text-success font-bold shadow-inner"
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
                        Fasilitas Pondok
                    </h2>
                </div>

                <div class="flex flex-wrap gap-3">
                    {#each detail.fasilitas_pondok as fasilitas (fasilitas)}
                        <span
                            class="badge font-medium px-4 py-3 rounded-xl bg-primary/10 text-primary border-primary/20"
                        >
                            {fasilitas}
                        </span>
                    {/each}
                </div>
            </div>
        </div>
    {/if}
</div>
