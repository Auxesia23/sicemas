<script lang="ts">
    import type { PageData } from "./$types";
    import type { ActivityResponse } from "$lib/schemas/activity";
    import { goto } from "$app/navigation";
    import Pagination from "$lib/components/ui/Pagination.svelte";

    let { data }: { data: PageData } = $props();
    let activities = $derived(data.activities);

    let selectedDate = $derived(data.selectedDate);

    let currentPage = $state(1);
    let itemsPerPage = $state(10);
    let searchTerm = $state("");

    $effect(() => {
        selectedDate = data.selectedDate;
    });

    $effect(() => {
        searchTerm;
        itemsPerPage;
        currentPage = 1;
    });

    const handleDateChange = () => {
        if (selectedDate) {
            goto(`?date=${selectedDate}`, { keepFocus: true });
        }
    };

    let filteredActivities = $derived(
        activities.filter((activity: ActivityResponse) => {
            const searchLower = searchTerm.toLowerCase();
            return (
                activity.user.toLowerCase().includes(searchLower) ||
                activity.action.toLowerCase().includes(searchLower) ||
                activity.target.toLowerCase().includes(searchLower) ||
                activity.entity.toLowerCase().includes(searchLower)
            );
        }),
    );

    let paginatedActivities = $derived(
        filteredActivities.slice(
            (currentPage - 1) * itemsPerPage,
            currentPage * itemsPerPage,
        ),
    );

    const formatTime = (isoString: string) => {
        const d = new Date(isoString);
        return d.toLocaleTimeString("id-ID");
    };
</script>

<svelte:head>
    <title>Log Aktivitas | SiCemas</title>
</svelte:head>

<div class="w-full">
    <div
        class="mb-8 flex flex-col gap-6 md:flex-row md:items-end md:justify-between"
    >
        <div>
            <div class="flex items-center gap-3 mb-2">
                <div
                    class="w-1.5 h-8 rounded-full bg-linear-to-b from-primary to-secondary"
                ></div>
                <h1
                    class="text-3xl font-extrabold tracking-tight text-base-content"
                >
                    Log Aktivitas
                </h1>
            </div>
            <p class="ml-5 font-medium text-base-content/60 text-sm">
                Rekam jejak aktivitas pengguna di sistem SiCemas
            </p>
        </div>
    </div>

    <div
        class="mb-6 rounded-2xl bg-base-100 p-4 shadow-sm border border-base-200 flex flex-col gap-4 sm:flex-row sm:items-center"
    >
        <div class="relative w-full sm:max-w-xs">
            <div class="text-xs font-bold text-base-content/70 mb-1 ml-1">
                Pilih Tanggal:
            </div>
            <input
                type="date"
                class="input input-bordered w-full rounded-xl bg-base-200/50 focus:bg-base-100 transition-colors font-medium"
                bind:value={selectedDate}
                onchange={handleDateChange}
            />
        </div>

        <div class="relative w-full sm:max-w-md self-end sm:ml-4">
            <div
                class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-4"
            >
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-5 w-5 text-base-content/40"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                    />
                </svg>
            </div>
            <input
                type="text"
                placeholder="Cari pengguna, aksi, atau target..."
                class="input input-bordered w-full rounded-xl bg-base-200/50 pl-11 focus:bg-base-100 transition-colors"
                bind:value={searchTerm}
            />
        </div>

        <div class="flex-1"></div>
        <div class="text-sm font-medium text-base-content/50 self-end mb-3">
            Total {filteredActivities.length} Aktivitas
        </div>
    </div>

    <div
        class="overflow-hidden rounded-2xl border border-base-200 bg-base-100 shadow-xl shadow-base-200/40"
    >
        {#if filteredActivities.length === 0}
            <div
                class="flex flex-col items-center justify-center py-16 px-4 text-center"
            >
                <div
                    class="mb-4 flex h-20 w-20 items-center justify-center rounded-full bg-base-200/50"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="h-10 w-10 text-base-content/30"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                        />
                    </svg>
                </div>
                <h3 class="text-lg font-bold text-base-content">
                    Tidak ada aktivitas di tanggal ini
                </h3>
                <p class="mt-1 text-sm text-base-content/60">
                    Coba pilih tanggal lain atau hapus filter pencarian.
                </p>
            </div>
        {:else}
            <div class="overflow-x-auto">
                <table class="table w-full">
                    <thead class="bg-base-200/50 text-base-content/70">
                        <tr
                            class="text-xs uppercase tracking-wider font-bold border-b-2 border-base-200"
                        >
                            <th class="py-4">Waktu</th>
                            <th>Pengguna</th>
                            <th>Aksi</th>
                            <th>Target</th>
                            <th>Entitas</th>
                        </tr>
                    </thead>
                    <tbody class="divide-y divide-base-200">
                        {#each paginatedActivities as log (log.id)}
                            <tr class="transition-colors hover:bg-base-200/30">
                                <td class="py-3 whitespace-nowrap">
                                    <div
                                        class="font-mono text-sm font-semibold text-base-content/80"
                                    >
                                        {formatTime(log.created_at)}
                                    </div>
                                </td>
                                <td>
                                    <div class="font-bold text-base-content">
                                        {log.user}
                                    </div>
                                </td>
                                <td>
                                    <span
                                        class="font-medium text-base-content/90"
                                    >
                                        {log.action}
                                    </span>
                                </td>
                                <td>
                                    <div class="text-sm text-base-content/80">
                                        {log.target}
                                    </div>
                                </td>
                                <td>
                                    <span
                                        class="badge badge-sm badge-outline text-xs font-bold border-base-content/20 text-base-content/70"
                                    >
                                        {log.entity}
                                    </span>
                                </td>
                            </tr>
                        {/each}
                    </tbody>
                </table>
            </div>

            <div class="border-t border-base-200 bg-base-100/50 p-4">
                <Pagination
                    bind:currentPage
                    bind:itemsPerPage
                    totalItems={filteredActivities.length}
                />
            </div>
        {/if}
    </div>
</div>
