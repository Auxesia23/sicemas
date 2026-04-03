<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { navigating, page } from "$app/state";
    import type { PageData } from "./$types";
    import type { SitusPublik } from "$lib/schemas/sites";

    let { data }: { data: PageData } = $props();

    let sites = $derived<SitusPublik[]>(data.sites || []);

    let searchTerm = $state("");
    let isLocating = $state(false);
    let userLocation = $state<{ lat: number; lng: number } | null>(null);

    let isLoading = $derived(navigating !== null || isLocating);

    let filteredSites = $derived(
        sites.filter((site) => {
            const searchLower = searchTerm.toLowerCase();
            return (
                site.nama.toLowerCase().includes(searchLower) ||
                site.jenis.toLowerCase().includes(searchLower) ||
                site.desa.toLowerCase().includes(searchLower)
            );
        }),
    );

    let currentPage = $state(1);
    const itemsPerPage = 9;

    $effect(() => {
        searchTerm;
        currentPage = 1;
    });

    let totalPages = $derived(Math.ceil(filteredSites.length / itemsPerPage));

    let paginatedSites = $derived(
        filteredSites.slice(
            (currentPage - 1) * itemsPerPage,
            currentPage * itemsPerPage,
        ),
    );

    function goToPage(pageNumber: number) {
        if (pageNumber >= 1 && pageNumber <= totalPages) {
            currentPage = pageNumber;
            document
                .getElementById("grid-situs")
                ?.scrollIntoView({ behavior: "smooth", block: "start" });
        }
    }

    function formatDistance(meters: number | null | undefined) {
        if (meters === null || meters === undefined) return null;
        if (meters < 1000) {
            return `${Math.round(meters)}m`;
        }
        return `${(meters / 1000).toFixed(1)}km`;
    }

    // Stats yang sudah dihitung
    let stats = $derived([
        { label: "Total Situs", value: sites.length, color: "primary" },
        {
            label: "Jenis Situs",
            value: new Set(sites.map((s) => s.jenis)).size,
            color: "secondary",
        },
        {
            label: "Dengan Foto",
            value: sites.filter((s) => s.foto_utama).length,
            color: "accent",
        },
        { label: "Terverifikasi", value: "100%", color: "primary" },
    ]);

    onMount(async () => {
        const existingLat = page.url.searchParams.get("lat");
        const existingLng = page.url.searchParams.get("lng");

        if (existingLat && existingLng) {
            userLocation = {
                lat: parseFloat(existingLat),
                lng: parseFloat(existingLng),
            };
            return;
        }

        isLocating = true;

        try {
            const position = await new Promise<GeolocationPosition>(
                (resolve, reject) => {
                    const timeoutId = setTimeout(
                        () => reject(new Error("GPS timeout")),
                        10000,
                    );

                    navigator.geolocation.getCurrentPosition(
                        (pos) => {
                            clearTimeout(timeoutId);
                            resolve(pos);
                        },
                        (err) => {
                            clearTimeout(timeoutId);
                            reject(err);
                        },
                        {
                            enableHighAccuracy: true,
                            timeout: 10000,
                            maximumAge: 0,
                        },
                    );
                },
            );

            const lat = position.coords.latitude;
            const lng = position.coords.longitude;
            userLocation = { lat, lng };

            await goto(`?lat=${lat}&lng=${lng}`, {
                replaceState: true,
                noScroll: true,
            });
        } catch (error: any) {
            console.warn("GPS not available:", error.message);
            userLocation = null;
        } finally {
            isLocating = false;
        }
    });
</script>

<svelte:head>
    <title>Direktori Situs Keagamaan - KUA Ciemas | SICEMAS</title>
    <meta
        name="description"
        content="Cari dan temukan daftar lengkap masjid, mushola, pondok pesantren, dan majelis taklim di wilayah Kecamatan Ciemas, Sukabumi. Data lengkap dengan lokasi terdekat."
    />
    <meta
        name="keywords"
        content="direktori masjid ciemas, cari masjid sukabumi, lokasi mushola ciemas, data pondok pesantren, situs keagamaan terdekat"
    />

    <meta
        property="og:title"
        content="Direktori Situs Keagamaan - KUA Ciemas"
    />
    <meta
        property="og:description"
        content="Temukan masjid, mushola, dan pondok pesantren di wilayah Ciemas dengan data lokasi terdekat."
    />
</svelte:head>

<div class="min-h-screen bg-base-100 relative">
    <div class="fixed inset-0 z-0 pointer-events-none">
        <div
            class="absolute inset-0 bg-linear-to-b from-base-100 to-base-200/50"
        ></div>

        <div
            class="absolute top-0 left-0 w-full h-150] bg-linear-to-b from-primary/15 via-primary/5 to-transparent"
        ></div>

        <div
            class="absolute top-[-10%] right-[-5%] w-96 h-96 bg-secondary/10 rounded-full blur-3xl"
        ></div>
        <div
            class="absolute top-[10%] left-[-5%] w-72 h-72 bg-accent/10 rounded-full blur-3xl"
        ></div>
    </div>

    <main
        class="relative z-10 container mx-auto px-4 pt-32 md:pt-44 pb-16 md:pb-24"
    >
        <section class="mb-16 text-center relative z-10">
            <div
                class="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-primary/10 text-primary text-sm font-medium mb-6 border border-primary/20"
            >
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-4 w-4"
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
                Direktori Lengkap
            </div>

            <h1
                class="text-4xl md:text-5xl lg:text-6xl font-bold text-base-content mb-6 leading-tight"
            >
                Jelajahi Situs Keagamaan
            </h1>

            <p
                class="mx-auto mb-12 max-w-2xl text-lg md:text-xl text-base-content/70 leading-relaxed"
            >
                Temukan dan jelajahi masjid, mushola, pondok pesantren, dan
                majelis taklim di wilayah Kecamatan Ciemas, Kabupaten Sukabumi
            </p>

            <div class="mx-auto max-w-2xl relative">
                <div class="relative group">
                    <div
                        class="absolute -inset-1 bg-linear-to-r from-primary to-secondary rounded-2xl blur opacity-25 group-hover:opacity-40 transition duration-1000 group-hover:duration-200"
                    ></div>
                    <div class="relative bg-base-100 rounded-2xl shadow-xl">
                        <div class="flex items-center px-6 py-4">
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                class="h-6 w-6 text-primary shrink-0"
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
                            <input
                                type="text"
                                class="w-full bg-transparent border-none focus:ring-0 text-lg px-4 placeholder-base-content/40 text-base-content"
                                placeholder="Cari nama situs, jenis, atau desa..."
                                bind:value={searchTerm}
                            />
                            {#if searchTerm}
                                <button
                                    class="p-2 hover:bg-base-200 rounded-full transition-colors"
                                    onclick={() => (searchTerm = "")}
                                    aria-label="Clear search"
                                >
                                    <svg
                                        xmlns="http://www.w3.org/2000/svg"
                                        class="h-5 w-5 text-base-content/50"
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
                            {/if}
                        </div>
                    </div>
                </div>

                {#if !searchTerm}
                    <div class="flex flex-wrap justify-center gap-2 mt-4">
                        {#each ["Masjid", "Mushola", "Ponpes", "MT"] as filter}
                            <button
                                class="px-4 py-1.5 rounded-full bg-base-100 border border-base-300 text-sm text-base-content/70 hover:border-primary/30 hover:text-primary transition-colors"
                                onclick={() => (searchTerm = filter)}
                            >
                                {filter}
                            </button>
                        {/each}
                    </div>
                {/if}
            </div>
        </section>

        {#if isLoading && sites.length === 0}
            <div class="flex items-center justify-center py-32">
                <div class="flex flex-col items-center gap-6">
                    <div class="relative">
                        <div
                            class="w-16 h-16 border-4 border-primary/20 border-t-primary rounded-full animate-spin"
                        ></div>
                        <div
                            class="absolute inset-0 w-16 h-16 border-4 border-secondary/20 border-b-secondary rounded-full animate-spin"
                            style="animation-direction: reverse; animation-duration: 1.5s;"
                        ></div>
                    </div>
                    <p class="text-base-content/70 text-lg animate-pulse">
                        {isLocating
                            ? "Mendeteksi lokasi Anda..."
                            : "Memuat data situs..."}
                    </p>
                </div>
            </div>
        {:else if sites.length === 0}
            <div class="flex items-center justify-center py-24">
                <div
                    class="flex flex-col items-center gap-6 text-center max-w-md mx-auto"
                >
                    <div
                        class="w-24 h-24 rounded-full bg-base-200 flex items-center justify-center"
                    >
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="h-12 w-12 text-base-content/30"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="1.5"
                                d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
                            />
                        </svg>
                    </div>
                    <div>
                        <h3 class="text-2xl font-bold text-base-content mb-2">
                            Belum Ada Data
                        </h3>
                        <p class="text-base-content/60">
                            Database situs keagamaan masih dalam proses
                            pengumpulan data. Silakan kembali lagi nanti.
                        </p>
                    </div>
                </div>
            </div>
        {:else if filteredSites.length === 0}
            <div class="flex items-center justify-center py-24">
                <div
                    class="flex flex-col items-center gap-6 text-center max-w-md mx-auto"
                >
                    <div
                        class="w-24 h-24 rounded-full bg-base-200 flex items-center justify-center"
                    >
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="h-12 w-12 text-base-content/30"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="1.5"
                                d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                            />
                        </svg>
                    </div>
                    <div>
                        <h3 class="text-2xl font-bold text-base-content mb-2">
                            Tidak Ditemukan
                        </h3>
                        <p class="text-base-content/60 mb-6">
                            Tidak ada situs yang cocok dengan pencarian "{searchTerm}"
                        </p>
                        <button
                            class="btn btn-primary btn-outline"
                            onclick={() => (searchTerm = "")}
                        >
                            Hapus Pencarian
                        </button>
                    </div>
                </div>
            </div>
        {:else}
            <div id="grid-situs" class="scroll-mt-24"></div>

            <div
                class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-8 px-2"
            >
                <p class="text-base-content/60 text-sm">
                    Menampilkan <span class="font-semibold text-base-content"
                        >{paginatedSites.length}</span
                    >
                    dari
                    <span class="font-semibold text-base-content"
                        >{filteredSites.length}</span
                    >
                    situs
                    {#if searchTerm}
                        untuk "<span class="font-semibold text-primary"
                            >{searchTerm}</span
                        >"
                    {/if}
                </p>
                {#if totalPages > 1}
                    <div
                        class="flex items-center gap-2 text-sm text-base-content/60 bg-base-100 px-4 py-2 rounded-full border border-base-200"
                    >
                        <span>Halaman</span>
                        <span class="font-bold text-primary">{currentPage}</span
                        >
                        <span>dari</span>
                        <span class="font-bold">{totalPages}</span>
                    </div>
                {/if}
            </div>

            <div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
                {#each paginatedSites as site (site.id)}
                    <article
                        class="group card bg-base-100 border border-base-200 shadow-lg hover:shadow-2xl transition-all duration-500 hover:-translate-y-2 overflow-hidden rounded-2xl"
                    >
                        <figure
                            class="relative h-56 w-full overflow-hidden bg-base-200"
                        >
                            {#if site.foto_utama}
                                <img
                                    src={site.foto_utama}
                                    alt={`Foto ${site.nama}`}
                                    class="h-full w-full object-cover transition-transform duration-700 group-hover:scale-110"
                                    loading="lazy"
                                />
                                <div
                                    class="absolute inset-0 bg-linear-to-t from-black/60 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300"
                                ></div>
                            {:else}
                                <div
                                    class="flex h-full w-full items-center justify-center bg-linear-to-br from-primary/10 to-secondary/10"
                                >
                                    <div class="text-center">
                                        <svg
                                            xmlns="http://www.w3.org/2000/svg"
                                            class="h-16 w-16 text-base-content/20 mx-auto mb-2"
                                            fill="none"
                                            viewBox="0 0 24 24"
                                            stroke="currentColor"
                                        >
                                            <path
                                                stroke-linecap="round"
                                                stroke-linejoin="round"
                                                stroke-width="1"
                                                d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
                                            />
                                        </svg>
                                        <span
                                            class="text-xs text-base-content/40"
                                            >Belum ada foto</span
                                        >
                                    </div>
                                </div>
                            {/if}

                            <div
                                class="absolute top-4 left-4 flex flex-col gap-2"
                            >
                                <span
                                    class="badge badge-primary shadow-lg text-white border-0"
                                >
                                    {site.jenis}
                                </span>
                                {#if site.jarak_meter !== null}
                                    <span
                                        class="badge badge-secondary shadow-lg text-white border-0 gap-1"
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
                                        {formatDistance(site.jarak_meter)}
                                    </span>
                                {/if}
                            </div>

                            <div
                                class="absolute bottom-4 right-4 translate-y-12 opacity-0 group-hover:translate-y-0 group-hover:opacity-100 transition-all duration-300"
                            >
                                <span
                                    class="btn btn-circle btn-sm btn-primary shadow-lg"
                                >
                                    <svg
                                        xmlns="http://www.w3.org/2000/svg"
                                        class="h-4 w-4"
                                        fill="none"
                                        viewBox="0 0 24 24"
                                        stroke="currentColor"
                                    >
                                        <path
                                            stroke-linecap="round"
                                            stroke-linejoin="round"
                                            stroke-width="2"
                                            d="M14 5l7 7m0 0l-7 7m7-7H3"
                                        />
                                    </svg>
                                </span>
                            </div>
                        </figure>

                        <div class="card-body p-5">
                            <div
                                class="flex items-start justify-between gap-2 mb-2"
                            >
                                <h3
                                    class="card-title text-lg font-bold text-base-content line-clamp-2 leading-tight group-hover:text-primary transition-colors"
                                >
                                    {site.nama}
                                </h3>
                            </div>

                            <div
                                class="flex items-start gap-2 text-sm text-base-content/60 mb-4"
                            >
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    class="mt-0.5 h-4 w-4 shrink-0 text-primary/70"
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
                                <span class="line-clamp-1"
                                    >{site.desa}, Kecamatan Ciemas</span
                                >
                            </div>

                            <a
                                href="/situs/{site.id}"
                                class="btn btn-block btn-outline btn-primary hover:btn-primary transition-all duration-300 group/btn"
                            >
                                <span>Lihat Detail</span>
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    class="h-4 w-4 transform group-hover/btn:translate-x-1 transition-transform"
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
                    </article>
                {/each}
            </div>

            {#if totalPages > 1}
                <div class="mt-16 flex justify-center">
                    <div
                        class="join bg-base-100 p-1 rounded-2xl shadow-lg border border-base-200"
                    >
                        <button
                            class="join-item btn btn-md px-6 hover:bg-base-200 disabled:opacity-50 disabled:cursor-not-allowed rounded-xl"
                            disabled={currentPage === 1}
                            onclick={() => goToPage(currentPage - 1)}
                        >
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                class="h-5 w-5 mr-2"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M15 19l-7-7 7-7"
                                />
                            </svg>
                            Sebelumnya
                        </button>

                        <div
                            class="join-item px-6 flex items-center bg-base-200/50 rounded-xl mx-1"
                        >
                            <span class="text-base-content/60 mr-2"
                                >Halaman</span
                            >
                            <span class="font-bold text-primary text-lg"
                                >{currentPage}</span
                            >
                            <span class="text-base-content/40 mx-2">/</span>
                            <span class="text-base-content/60"
                                >{totalPages}</span
                            >
                        </div>

                        <button
                            class="join-item btn btn-md px-6 hover:bg-base-200 disabled:opacity-50 disabled:cursor-not-allowed rounded-xl"
                            disabled={currentPage === totalPages}
                            onclick={() => goToPage(currentPage + 1)}
                        >
                            Selanjutnya
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                class="h-5 w-5 ml-2"
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
                        </button>
                    </div>
                </div>
            {/if}

            <section class="mt-24 relative">
                <div
                    class="absolute inset-0 bg-linear-to-r from-primary/5 via-secondary/5 to-accent/5 rounded-3xl -z-10"
                    aria-hidden="true"
                ></div>
                <div
                    class="rounded-3xl bg-base-100/80 backdrop-blur-sm p-8 md:p-12 shadow-xl border border-base-200"
                >
                    <div class="grid grid-cols-2 md:grid-cols-4 gap-8 md:gap-4">
                        {#each stats as stat}
                            <div class="text-center group">
                                <div class="relative inline-block mb-3">
                                    <div
                                        class="absolute inset-0 bg-{stat.color}/20 rounded-full blur-xl group-hover:blur-2xl transition-all duration-300"
                                    ></div>
                                    <div
                                        class="relative text-4xl md:text-5xl font-bold text-{stat.color} font-display"
                                    >
                                        {stat.value}
                                    </div>
                                </div>
                                <div
                                    class="text-sm md:text-base text-base-content/60 font-medium"
                                >
                                    {stat.label}
                                </div>
                            </div>
                        {/each}
                    </div>
                </div>
            </section>
        {/if}
    </main>
</div>

<style>
    .line-clamp-2 {
        display: -webkit-box;
        line-clamp: 2;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
    }

    .line-clamp-1 {
        display: -webkit-box;
        line-clamp: 1;
        -webkit-line-clamp: 1;
        -webkit-box-orient: vertical;
        overflow: hidden;
    }

    .font-display {
        font-feature-settings: "tnum";
        font-variant-numeric: tabular-nums;
    }

    @keyframes spin-slow {
        to {
            transform: rotate(360deg);
        }
    }

    .card {
        transform-style: preserve-3d;
    }
</style>
