<script lang="ts">
    import type { PageData } from "./$types";

    let { data }: { data: PageData } = $props();

    let dashboard = $derived(data.dashboard);

    let statsConfig = $derived([
        {
            title: "Total Situs",
            value: dashboard.stats.total_situs,
            icon: "M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4",
            color: "primary",
        },
        {
            title: "Jenis Situs",
            value: dashboard.stats.total_jenis,
            icon: "M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z",
            color: "secondary",
        },
        {
            title: "Petugas Aktif",
            value: dashboard.stats.petugas_aktif,
            icon: "M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z",
            color: "accent",
        },
        {
            title: "Menunggu Verifikasi",
            value: dashboard.stats.menunggu_verifikasi,
            icon: "M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z",
            color: "warning",
        },
    ]);

    function timeAgo(dateString: string | null): string {
        if (!dateString) return "-";
        const date = new Date(dateString);
        const now = new Date();
        const seconds = Math.round((now.getTime() - date.getTime()) / 1000);
        const minutes = Math.round(seconds / 60);
        const hours = Math.round(minutes / 60);
        const days = Math.round(hours / 24);

        if (seconds < 60) return "Baru saja";
        if (minutes < 60) return `${minutes} menit yang lalu`;
        if (hours < 24) return `${hours} jam yang lalu`;
        if (days === 1) return "Kemarin";
        return `${days} hari yang lalu`;
    }

    function formatDate(dateString: string | null): string {
        if (!dateString) return "-";
        const date = new Date(dateString);
        return date.toLocaleDateString("id-ID", {
            day: "numeric",
            month: "short",
            year: "numeric",
        });
    }

    function getStatusBadgeClass(status: string | null): string {
        switch (status?.toLowerCase()) {
            case "menunggu":
                return "badge-warning";
            case "terverifikasi":
                return "badge-success";
            case "ditolak":
                return "badge-error";
            default:
                return "badge-ghost";
        }
    }
</script>

<svelte:head>
    <title>Dashboard | SiCemas</title>
    <meta name="description" content="Sistem Informasi Keagamaan KUA Ciemas" />
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
                    Dashboard Admin
                </h1>
            </div>
            <p class="ml-5 font-medium text-base-content/60 text-sm">
                Selamat datang di sistem pendataan situs keagamaan KUA Kecamatan
                Ciemas
            </p>
        </div>
    </div>

    <div class="mb-8 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
        {#each statsConfig as stat (stat.title)}
            <div
                class="card rounded-2xl border border-base-200 bg-base-100 shadow-xl shadow-base-200/40 transition-all duration-300 hover:-translate-y-1 hover:shadow-2xl hover:shadow-base-200/50"
            >
                <div class="card-body p-6">
                    <div class="flex items-center justify-between">
                        <div>
                            <p
                                class="text-xs font-bold uppercase tracking-wider text-base-content/50"
                            >
                                {stat.title}
                            </p>
                            <h3
                                class="mt-1 text-3xl font-extrabold text-{stat.color}"
                            >
                                {stat.value}
                            </h3>
                        </div>
                        <div
                            class={`flex h-14 w-14 items-center justify-center rounded-2xl bg-${stat.color}/10 text-${stat.color} shadow-sm`}
                        >
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                class="h-7 w-7"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d={stat.icon}
                                />
                            </svg>
                        </div>
                    </div>
                </div>
            </div>
        {/each}
    </div>

    <div class="mb-8 grid grid-cols-1 gap-6 lg:grid-cols-3">
        <div
            class="card rounded-2xl border border-base-200 bg-base-100 shadow-xl shadow-base-200/40 lg:col-span-2"
        >
            <div class="card-body p-6 sm:p-8">
                <h2
                    class="card-title text-xl font-extrabold text-base-content border-b border-base-200 pb-4"
                >
                    Statistik Berdasarkan Jenis
                </h2>

                <div class="mt-4 flex flex-col justify-center space-y-6">
                    {#if dashboard.statistik_jenis.length > 0}
                        {@const maxSitus = Math.max(
                            ...dashboard.statistik_jenis.map(
                                (d) => d.jumlah_situs,
                            ),
                        )}

                        {#each dashboard.statistik_jenis as item}
                            <div class="flex flex-col gap-2">
                                <div class="flex justify-between text-sm">
                                    <span
                                        class="font-bold text-base-content/80"
                                    >
                                        {item.nama}
                                    </span>
                                    <span class="font-bold text-base-content">
                                        {item.jumlah_situs}
                                        <span
                                            class="font-normal text-base-content/60"
                                            >Situs</span
                                        >
                                    </span>
                                </div>
                                <div
                                    class="h-3 w-full overflow-hidden rounded-full bg-base-200/50"
                                >
                                    <div
                                        class="h-full rounded-full bg-linear-to-r from-primary to-secondary transition-all duration-1000 ease-out shadow-sm"
                                        style="width: {maxSitus === 0
                                            ? 0
                                            : (item.jumlah_situs / maxSitus) *
                                              100}%"
                                    ></div>
                                </div>
                            </div>
                        {/each}
                    {:else}
                        <div
                            class="flex flex-col items-center justify-center py-10 opacity-60"
                        >
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                class="h-12 w-12 mb-3"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
                                />
                            </svg>
                            <span class="text-sm font-medium"
                                >Belum ada data statistik jenis situs.</span
                            >
                        </div>
                    {/if}
                </div>
            </div>
        </div>

        <div
            class="card rounded-2xl border border-base-200 bg-base-100 shadow-xl shadow-base-200/40"
        >
            <div class="card-body p-6 sm:p-8">
                <h2
                    class="card-title text-xl font-extrabold text-base-content border-b border-base-200 pb-4"
                >
                    Aktivitas Terbaru
                </h2>

                <div class="mt-4 space-y-5">
                    {#each dashboard.recent_activities as activity (activity.id)}
                        <div class="flex items-start space-x-4">
                            <div class="avatar placeholder">
                                <div
                                    class="flex h-10 w-10 items-center justify-center rounded-full bg-linear-to-br from-primary/20 to-secondary/20 text-primary font-bold shadow-sm"
                                >
                                    {activity.user.charAt(0).toUpperCase()}
                                </div>
                            </div>
                            <div class="min-w-0 flex-1 pt-0.5">
                                <p
                                    class="text-sm leading-tight text-base-content/80"
                                >
                                    <span class="font-bold text-base-content"
                                        >{activity.user}</span
                                    >
                                    {activity.action}
                                    <span class="font-bold text-base-content"
                                        >{activity.target}</span
                                    >
                                </p>
                                <p
                                    class="mt-1 text-xs font-medium text-base-content/50 flex items-center gap-1"
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
                                            d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                                        />
                                    </svg>
                                    {timeAgo(activity.created_at)}
                                </p>
                            </div>
                        </div>
                    {:else}
                        <div
                            class="text-center py-10 opacity-60 flex flex-col items-center"
                        >
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                class="h-10 w-10 mb-2"
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
                            <span class="text-sm font-medium"
                                >Belum ada aktivitas terbaru.</span
                            >
                        </div>
                    {/each}
                </div>
            </div>
        </div>
    </div>

    <div
        class="card overflow-hidden rounded-2xl border border-base-200 bg-base-100 shadow-xl shadow-base-200/40"
    >
        <div class="card-body p-0 sm:p-2">
            <div
                class="flex items-center justify-between p-6 border-b border-base-200"
            >
                <h2 class="card-title text-xl font-extrabold text-base-content">
                    Situs Baru Ditambahkan
                </h2>
                <a
                    href="/dashboard/situs"
                    class="btn btn-sm btn-ghost text-primary hover:bg-primary/10 rounded-lg"
                >
                    Lihat Semua
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="h-4 w-4 ml-1"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M9 5l7 7-7 7"
                        />
                    </svg>
                </a>
            </div>

            <div class="overflow-x-auto">
                <table class="table w-full">
                    <thead class="bg-base-200/30 text-base-content/70">
                        <tr class="text-xs uppercase tracking-wider font-bold">
                            <th class="py-4 pl-6">Situs & Lokasi</th>
                            <th>Kategori</th>
                            <th>Pendata</th>
                            <th class="text-center">Status</th>
                        </tr>
                    </thead>
                    <tbody class="divide-y divide-base-200">
                        {#each dashboard.recent_sites as site (site.id)}
                            <tr class="transition-colors hover:bg-base-200/30">
                                <td class="py-4 pl-6">
                                    <div
                                        class="font-extrabold text-base-content text-sm mb-0.5"
                                    >
                                        {site.nama}
                                    </div>
                                    <div
                                        class="text-xs text-base-content/60 flex items-center gap-1"
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
                                                d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"
                                            />
                                            <path
                                                stroke-linecap="round"
                                                stroke-linejoin="round"
                                                stroke-width="2"
                                                d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"
                                            />
                                        </svg>
                                        <span
                                            class="truncate max-w-50"
                                            title={site.lokasi}
                                            >{site.lokasi}</span
                                        >
                                    </div>
                                </td>
                                <td>
                                    <span
                                        class="badge border-primary/20 bg-primary/10 text-primary text-xs font-semibold px-2.5 py-3"
                                    >
                                        {site.jenis}
                                    </span>
                                </td>
                                <td>
                                    <div
                                        class="text-sm font-medium text-base-content/80"
                                    >
                                        {site.pendata}
                                    </div>
                                    <div
                                        class="text-[10px] text-base-content/40 mt-0.5"
                                    >
                                        {formatDate(site.updated_at)}
                                    </div>
                                </td>
                                <td class="text-center pr-6">
                                    <div
                                        class={`badge badge-sm font-medium px-2.5 py-3 ${
                                            site.status === "terverifikasi"
                                                ? "bg-success/10 text-success border-success/20"
                                                : site.status === "menunggu"
                                                  ? "bg-warning/10 text-warning border-warning/20"
                                                  : "bg-error/10 text-error border-error/20"
                                        }`}
                                    >
                                        <div
                                            class={`w-1.5 h-1.5 rounded-full mr-1.5 ${
                                                site.status === "terverifikasi"
                                                    ? "bg-success"
                                                    : site.status === "menunggu"
                                                      ? "bg-warning"
                                                      : "bg-error"
                                            }`}
                                        ></div>
                                        {site.status.charAt(0).toUpperCase() +
                                            site.status.slice(1)}
                                    </div>
                                </td>
                            </tr>
                        {:else}
                            <tr>
                                <td
                                    colspan="4"
                                    class="text-center py-12 text-base-content/50"
                                >
                                    <div
                                        class="flex flex-col items-center justify-center"
                                    >
                                        <svg
                                            xmlns="http://www.w3.org/2000/svg"
                                            class="h-10 w-10 mb-2 opacity-50"
                                            fill="none"
                                            viewBox="0 0 24 24"
                                            stroke="currentColor"
                                        >
                                            <path
                                                stroke-linecap="round"
                                                stroke-linejoin="round"
                                                stroke-width="2"
                                                d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                                            />
                                        </svg>
                                        <span
                                            >Belum ada data situs yang
                                            terdaftar.</span
                                        >
                                    </div>
                                </td>
                            </tr>
                        {/each}
                    </tbody>
                </table>
            </div>
        </div>
    </div>
</div>
