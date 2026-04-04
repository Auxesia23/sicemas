<script lang="ts">
    import { onMount } from "svelte";
    import "leaflet/dist/leaflet.css";
    import type { PageData } from "./$types";

    let { data }: { data: PageData } = $props();
    let site = $derived(data.site);

    let mapContainer: HTMLElement | undefined = $state();
    let map: any;
    let marker: any;
    let selectedImage = $state<string | null>(null);

    const facilityIcons: Record<string, string> = {
        tempat_wudhu: `<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 14.66V20a2 2 0 01-2 2H4a2 2 0 01-2-2V6a2 2 0 012-2h2.34M20 14.66a3 3 0 01-3-3V6a3 3 0 013-3h.01M20 14.66a3 3 0 003-3V6a3 3 0 00-3-3h-.01M10 14v6m-4-3h8"/></svg>`,
        musholla: `<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"/></svg>`,
        perpustakaan: `<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/></svg>`,
        aula: `<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"/></svg>`,
        kantor: `<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 13.255A23.931 23.931 0 0112 15c-3.183 0-6.22-.62-9-1.745M16 6V4a2 2 0 00-2-2h-4a2 2 0 00-2 2v2m4 6h.01M5 20h14a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/></svg>`,
        parkir: `<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 10l7-7m0 0l7 7m-7-7v18"/></svg>`,
        default: `<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/></svg>`,
    };

    function getFacilityIcon(key: string): string {
        return facilityIcons[key.toLowerCase()] || facilityIcons["default"];
    }

    function formatKey(key: string): string {
        return key
            .split("_")
            .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
            .join(" ");
    }

    function openLightbox(image: string) {
        selectedImage = image;
        document.body.style.overflow = "hidden";
    }

    function closeLightbox() {
        selectedImage = null;
        document.body.style.overflow = "";
    }

    onMount(async () => {
        if (
            !site ||
            site.latitude === null ||
            site.longitude === null ||
            !mapContainer
        )
            return;

        const L = await import("leaflet");

        delete (L.Icon.Default.prototype as any)._getIconUrl;
        L.Icon.Default.mergeOptions({
            iconRetinaUrl:
                "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-icon-2x.png",
            iconUrl:
                "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-icon.png",
            shadowUrl:
                "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-shadow.png",
        });

        map = L.map(mapContainer, {
            dragging: true,
            touchZoom: true,
            scrollWheelZoom: false,
            doubleClickZoom: false,
            zoomControl: true,
            attributionControl: true,
        }).setView([site.latitude, site.longitude], 15);

        L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
            attribution: "© OpenStreetMap contributors",
        }).addTo(map);

        marker = L.marker([site.latitude, site.longitude]).addTo(map);

        setTimeout(() => {
            if (map) map.invalidateSize();
        }, 200);
    });

    let isCopied = $state(false);

    async function handleShare() {
        const shareData = {
            title: site.nama,
            text: `Lihat ${site.nama} di SICEMAS`,
            url: window.location.href,
        };

        if (navigator.share) {
            try {
                await navigator.share(shareData);
            } catch (err: any) {
                console.log(err);
                if (err.name !== "AbortError") {
                    await fallbackCopy();
                }
            }
        } else {
            await fallbackCopy();
        }
    }

    async function fallbackCopy() {
        try {
            await navigator.clipboard.writeText(window.location.href);
            isCopied = true;
            setTimeout(() => {
                isCopied = false;
            }, 2000);
        } catch (err) {
            console.log("Gagal copy link: ", err);
            alert("Gagal menyalin link.");
        }
    }
</script>

<svelte:head>
    <title
        >{site?.nama
            ? `${site.nama} - Direktori Situs Keagamaan`
            : "Detail Situs"} | SICEMAS KUA Ciemas</title
    >
    {#if site?.alamat_lengkap}
        <meta
            name="description"
            content={`${site.nama} - ${site.jenis} di ${site.desa}, Kecamatan Ciemas. ${site.alamat_lengkap}. Daya tampung ${site.daya_tampung} orang, luas tanah ${site.luas_tanah}m².`}
        />
    {:else}
        <meta
            name="description"
            content="Detail informasi situs keagamaan di Kecamatan Ciemas, Kabupaten Sukabumi."
        />
    {/if}
    <meta
        name="keywords"
        content={`${site?.nama || ""}, ${site?.jenis || ""}, ${site?.desa || ""}, masjid ciemas, situs keagamaan sukabumi`}
    />

    <meta property="og:title" content={site?.nama} />
    <meta property="og:description" content={site?.alamat_lengkap} />
    <meta property="og:type" content="article" />
    <meta
        property="og:image"
        content={site?.galeri?.[0] || "/default-og.jpg"}
    />

    <link
        rel="canonical"
        href={`https://sicemas.ciemas.desa.id/situs/${site?.id || ""}`}
    />
</svelte:head>

{#if selectedImage}
    <div
        class="fixed inset-0 z-50 bg-black/95 flex items-center justify-center p-4 backdrop-blur-sm transition-all"
        onclick={closeLightbox}
        onkeydown={(e) => e.key === "Escape" && closeLightbox()}
        role="button"
        tabindex="0"
        aria-label="Close lightbox"
    >
        <button
            class="absolute top-4 right-4 text-white/60 hover:text-white p-2 transition-colors bg-white/10 hover:bg-white/20 rounded-full"
            onclick={closeLightbox}
            aria-label="Close"
        >
            <svg
                class="w-6 h-6"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                ><path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M6 18L18 6M6 6l12 12"
                /></svg
            >
        </button>
        <img
            src={selectedImage}
            alt="Full size"
            class="max-w-full max-h-[90vh] object-contain rounded-2xl shadow-2xl ring-1 ring-white/10"
        />
    </div>
{/if}

<div class="min-h-screen bg-base-100 relative">
    <div class="fixed inset-0 z-0 pointer-events-none">
        <div
            class="absolute inset-0 bg-linear-to-b from-base-100 to-base-200/50"
        ></div>
        <div
            class="absolute top-0 left-0 w-full h-125 bg-linear-to-b from-primary/5 via-secondary/5 to-transparent"
        ></div>
        <div
            class="absolute top-20 right-[-5%] w-96 h-96 bg-secondary/10 rounded-full blur-[100px] opacity-60"
        ></div>
        <div
            class="absolute top-60 left-[-5%] w-72 h-72 bg-accent/10 rounded-full blur-[80px] opacity-60"
        ></div>
    </div>

    {#if site}
        <div class="relative z-10">
            <div class="container mx-auto px-4 pt-8 pb-4 sm:px-6 lg:px-8">
                <nav aria-label="Breadcrumb">
                    <ol
                        class=" items-center gap-2.5 text-sm font-medium text-base-content/60 flex-wrap bg-base-200/30 px-4 py-2 rounded-xl inline-flex border border-base-200 backdrop-blur-sm"
                    >
                        <li>
                            <a
                                href="/"
                                class="hover:text-primary transition-colors flex items-center gap-1.5"
                            >
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    class="h-4 w-4"
                                    fill="none"
                                    viewBox="0 0 24 24"
                                    stroke="currentColor"
                                    ><path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"
                                    /></svg
                                >
                                Beranda
                            </a>
                        </li>
                        <li>
                            <svg
                                class="w-4 h-4 opacity-40"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                                ><path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M9 5l7 7-7 7"
                                /></svg
                            >
                        </li>
                        <li>
                            <a
                                href="/situs"
                                class="hover:text-primary transition-colors"
                                >Direktori</a
                            >
                        </li>
                        <li>
                            <svg
                                class="w-4 h-4 opacity-40"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                                ><path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M9 5l7 7-7 7"
                                /></svg
                            >
                        </li>
                        <li
                            class="text-primary font-bold truncate max-w-50 sm:max-w-75"
                        >
                            {site.nama}
                        </li>
                    </ol>
                </nav>
            </div>

            <main class="container mx-auto px-4 py-6 sm:px-6 lg:px-8 mb-20">
                <section class="mb-12">
                    <div class="flex flex-col lg:flex-row gap-8 items-start">
                        <div class="w-full lg:w-2/3">
                            {#if site.galeri?.[0]}
                                <div
                                    class="relative group rounded-3xl overflow-hidden shadow-2xl shadow-base-300/50 border border-base-200 bg-base-200"
                                >
                                    <img
                                        src={site.galeri?.[0]}
                                        alt={site.nama}
                                        class="w-full h-72 md:h-112.5 object-cover transition-transform duration-700 group-hover:scale-105"
                                    />
                                    <div
                                        class="absolute inset-0 bg-linear-to-t from-black/60 via-black/10 to-transparent opacity-80"
                                    ></div>

                                    {#if site.galeri && site.galeri.length > 0}
                                        <button
                                            class="absolute bottom-5 right-5 btn btn-circle btn-primary border-none bg-primary/90 hover:bg-primary backdrop-blur-md shadow-xl transition-transform hover:-translate-y-1"
                                            onclick={() =>
                                                openLightbox(site.galeri?.[0])}
                                            aria-label="View full image"
                                        >
                                            <svg
                                                class="w-5 h-5 text-white"
                                                fill="none"
                                                stroke="currentColor"
                                                viewBox="0 0 24 24"
                                            >
                                                <path
                                                    stroke-linecap="round"
                                                    stroke-linejoin="round"
                                                    stroke-width="2"
                                                    d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0zM10 7v3m0 0v3m0-3h3m-3 0H7"
                                                />
                                            </svg>
                                        </button>
                                    {/if}

                                    <div class="absolute top-5 left-5">
                                        <span
                                            class="badge badge-lg bg-base-100/90 text-base-content font-bold border-0 shadow-lg backdrop-blur-md px-4 py-3"
                                        >
                                            {site.jenis}
                                        </span>
                                    </div>
                                </div>
                            {:else}
                                <div
                                    class="w-full h-72 md:h-112.5 flex flex-col items-center justify-center rounded-3xl border-2 border-dashed border-base-300 bg-base-200/30"
                                >
                                    <div
                                        class="w-20 h-20 rounded-full bg-base-300/50 flex items-center justify-center mb-4"
                                    >
                                        <svg
                                            class="w-10 h-10 text-base-content/30"
                                            fill="none"
                                            stroke="currentColor"
                                            viewBox="0 0 24 24"
                                        >
                                            <path
                                                stroke-linecap="round"
                                                stroke-linejoin="round"
                                                stroke-width="2"
                                                d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
                                            />
                                        </svg>
                                    </div>
                                    <p
                                        class="text-base-content/40 font-bold text-lg"
                                    >
                                        Belum ada foto
                                    </p>
                                </div>
                            {/if}
                        </div>

                        <div class="w-full lg:w-1/3">
                            <div
                                class="bg-base-100/80 backdrop-blur-xl rounded-3xl p-6 md:p-8 shadow-xl shadow-base-200/50 border border-base-200 lg:sticky lg:top-24 transition-all"
                            >
                                <h1
                                    class="text-3xl md:text-4xl font-extrabold text-base-content mb-4 leading-tight tracking-tight"
                                >
                                    {site.nama}
                                </h1>

                                <div
                                    class="flex items-start gap-3 text-base-content/70 mb-8 p-4 bg-base-200/50 rounded-2xl border border-base-200"
                                >
                                    <div
                                        class="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center shrink-0 mt-0.5"
                                    >
                                        <svg
                                            class="w-4 h-4 text-primary"
                                            fill="none"
                                            stroke="currentColor"
                                            viewBox="0 0 24 24"
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
                                    </div>
                                    <div class="flex flex-col">
                                        <span
                                            class="text-xs font-bold uppercase tracking-wider text-base-content/50 mb-1"
                                            >Alamat Lengkap</span
                                        >
                                        <span
                                            class="leading-relaxed font-medium"
                                            >{site.alamat_lengkap}</span
                                        >
                                    </div>
                                </div>

                                <div class="grid grid-cols-2 gap-4 mb-8">
                                    <div
                                        class="bg-primary/5 rounded-2xl p-5 border border-primary/10 transition-transform hover:-translate-y-1 hover:shadow-md"
                                    >
                                        <div
                                            class="flex items-center gap-2 text-primary mb-2"
                                        >
                                            <svg
                                                class="w-5 h-5"
                                                fill="none"
                                                stroke="currentColor"
                                                viewBox="0 0 24 24"
                                                ><path
                                                    stroke-linecap="round"
                                                    stroke-linejoin="round"
                                                    stroke-width="2"
                                                    d="M4 8V4m0 0h4M4 4l5 5m11-1V4m0 0h-4m4 0l-5 5M4 16v4m0 0h4m-4 0l5-5m11 5l-5-5m5 5v-4m0 4h-4"
                                                /></svg
                                            >
                                            <span
                                                class="text-xs font-bold uppercase tracking-wider"
                                                >Luas Tanah</span
                                            >
                                        </div>
                                        <div
                                            class="text-2xl font-extrabold text-base-content"
                                        >
                                            {site.luas_tanah}
                                            <span
                                                class="text-sm font-medium text-base-content/50"
                                                >m²</span
                                            >
                                        </div>
                                    </div>
                                    <div
                                        class="bg-secondary/5 rounded-2xl p-5 border border-secondary/10 transition-transform hover:-translate-y-1 hover:shadow-md"
                                    >
                                        <div
                                            class="flex items-center gap-2 text-secondary mb-2"
                                        >
                                            <svg
                                                class="w-5 h-5"
                                                fill="none"
                                                stroke="currentColor"
                                                viewBox="0 0 24 24"
                                                ><path
                                                    stroke-linecap="round"
                                                    stroke-linejoin="round"
                                                    stroke-width="2"
                                                    d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
                                                /></svg
                                            >
                                            <span
                                                class="text-xs font-bold uppercase tracking-wider"
                                                >Kapasitas</span
                                            >
                                        </div>
                                        <div
                                            class="text-2xl font-extrabold text-base-content"
                                        >
                                            {site.daya_tampung}
                                            <span
                                                class="text-sm font-medium text-base-content/50"
                                                >org</span
                                            >
                                        </div>
                                    </div>
                                </div>

                                <div class="space-y-3">
                                    {#if site.latitude && site.longitude}
                                        <a
                                            href={`https://www.google.com/maps/search/?api=1&query=${site.latitude},${site.longitude}`}
                                            target="_blank"
                                            rel="noopener noreferrer"
                                            class="btn btn-primary w-full gap-2 rounded-xl h-12 shadow-lg shadow-primary/20 transition-transform hover:-translate-y-0.5"
                                        >
                                            <svg
                                                class="w-5 h-5"
                                                fill="none"
                                                stroke="currentColor"
                                                viewBox="0 0 24 24"
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
                                            Buka di Google Maps
                                        </a>
                                    {/if}
                                    <div class="flex gap-3">
                                        {#if site.latitude && site.longitude}
                                            <a
                                                href={`https://www.google.com/maps/dir/?api=1&destination=${site.latitude},${site.longitude}`}
                                                target="_blank"
                                                rel="noopener noreferrer"
                                                class="btn bg-base-200/50 hover:bg-base-200 border-base-200 flex-1 rounded-xl h-12 text-base-content/80 font-bold"
                                            >
                                                Rute Jalan
                                            </a>
                                        {/if}
                                        <button
                                            class="btn bg-base-200/50 hover:bg-base-200 border-base-200 flex-1 rounded-xl h-12 text-base-content/80 font-bold gap-2"
                                            onclick={handleShare}
                                        >
                                            <svg
                                                class="w-4 h-4"
                                                fill="none"
                                                stroke="currentColor"
                                                viewBox="0 0 24 24"
                                                ><path
                                                    stroke-linecap="round"
                                                    stroke-linejoin="round"
                                                    stroke-width="2"
                                                    d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z"
                                                /></svg
                                            >
                                            {isCopied ? "Disalin!" : "Bagikan"}
                                        </button>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </section>

                {#if site.galeri && site.galeri.length > 1}
                    <section class="mb-16">
                        <div class="flex items-center justify-between mb-6">
                            <h2
                                class="text-2xl font-extrabold tracking-tight text-base-content flex items-center gap-3"
                            >
                                <div
                                    class="w-1.5 h-6 rounded-full bg-linear-to-b from-primary to-secondary"
                                ></div>
                                Galeri Foto
                            </h2>
                            <span
                                class="badge badge-lg bg-base-200/50 border-base-200 font-bold"
                                >{site.galeri.length} Foto</span
                            >
                        </div>

                        <div
                            class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 sm:gap-6"
                        >
                            {#each site.galeri as foto, index (foto)}
                                <button
                                    class="group relative aspect-square rounded-3xl overflow-hidden bg-base-200 shadow-md hover:shadow-xl border border-base-200 transition-all duration-300 focus:outline-none focus:ring-4 focus:ring-primary/30"
                                    onclick={() => openLightbox(foto)}
                                    aria-label={`View image ${index + 1}`}
                                >
                                    <img
                                        src={foto}
                                        alt={`${site.nama} - Foto ${index + 1}`}
                                        class="w-full h-full object-cover transition-transform duration-700 group-hover:scale-110"
                                        loading="lazy"
                                    />
                                    <div
                                        class="absolute inset-0 bg-linear-to-t from-black/60 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300"
                                    ></div>
                                    <div
                                        class="absolute inset-0 flex items-center justify-center pointer-events-none"
                                    >
                                        <div
                                            class="w-12 h-12 rounded-full bg-white/20 backdrop-blur-md flex items-center justify-center opacity-0 group-hover:opacity-100 transform scale-50 group-hover:scale-100 transition-all duration-300"
                                        >
                                            <svg
                                                class="w-5 h-5 text-white"
                                                fill="none"
                                                stroke="currentColor"
                                                viewBox="0 0 24 24"
                                            >
                                                <path
                                                    stroke-linecap="round"
                                                    stroke-linejoin="round"
                                                    stroke-width="2"
                                                    d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0zM10 7v3m0 0v3m0-3h3m-3 0H7"
                                                />
                                            </svg>
                                        </div>
                                    </div>
                                </button>
                            {/each}
                        </div>
                    </section>
                {/if}

                {#if site.fasilitas && Object.keys(site.fasilitas).length > 0}
                    <section class="mb-16">
                        <h2
                            class="text-2xl font-extrabold tracking-tight text-base-content flex items-center gap-3 mb-8"
                        >
                            <div
                                class="w-1.5 h-6 rounded-full bg-linear-to-b from-primary to-secondary"
                            ></div>
                            Fasilitas & Layanan
                        </h2>
                        <div
                            class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"
                        >
                            {#each Object.entries(site.fasilitas) as [key, values]}
                                {#if Array.isArray(values) && values.length > 0}
                                    <div
                                        class="bg-base-100/80 backdrop-blur-sm rounded-3xl p-6 shadow-lg shadow-base-200/30 border border-base-200 hover:-translate-y-1 transition-transform duration-300"
                                    >
                                        <div
                                            class="flex items-center gap-4 mb-5 border-b border-base-200 pb-4"
                                        >
                                            <div
                                                class="w-12 h-12 rounded-2xl bg-linear-to-br from-primary/10 to-secondary/10 flex items-center justify-center text-primary shadow-inner"
                                            >
                                                {@html getFacilityIcon(key)}
                                            </div>
                                            <h3
                                                class="text-lg font-extrabold text-base-content"
                                            >
                                                {formatKey(key)}
                                            </h3>
                                        </div>
                                        <div class="flex flex-wrap gap-2.5">
                                            {#each values as item (item)}
                                                <span
                                                    class="px-3.5 py-1.5 rounded-xl bg-base-200/50 text-sm font-medium text-base-content/80 border border-base-200"
                                                >
                                                    {item}
                                                </span>
                                            {/each}
                                        </div>
                                    </div>
                                {/if}
                            {/each}
                        </div>
                    </section>
                {/if}

                {#if site.latitude !== null && site.longitude !== null}
                    <section class="mb-16">
                        <h2
                            class="text-2xl font-extrabold tracking-tight text-base-content flex items-center gap-3 mb-8"
                        >
                            <div
                                class="w-1.5 h-6 rounded-full bg-linear-to-b from-primary to-secondary"
                            ></div>
                            Lokasi Presisi
                        </h2>
                        <div
                            class="bg-base-100 rounded-3xl overflow-hidden shadow-2xl shadow-base-200/40 border border-base-200 p-2"
                        >
                            <div
                                class="h-100 w-full relative rounded-2xl overflow-hidden"
                            >
                                <div
                                    bind:this={mapContainer}
                                    class="h-full w-full z-0"
                                ></div>

                                <div
                                    class="absolute bottom-4 left-4 right-4 md:left-auto md:right-4 md:w-80 bg-base-100/90 backdrop-blur-xl rounded-2xl p-5 shadow-xl border border-white/20 z-400"
                                >
                                    <div class="flex items-start gap-4">
                                        <div
                                            class="w-10 h-10 rounded-xl bg-primary/10 flex items-center justify-center shrink-0 text-primary mt-0.5"
                                        >
                                            <svg
                                                class="w-5 h-5"
                                                fill="none"
                                                stroke="currentColor"
                                                viewBox="0 0 24 24"
                                                ><path
                                                    stroke-linecap="round"
                                                    stroke-linejoin="round"
                                                    stroke-width="2"
                                                    d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"
                                                /><path
                                                    stroke-linecap="round"
                                                    stroke-linejoin="round"
                                                    stroke-width="2"
                                                    d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"
                                                /></svg
                                            >
                                        </div>
                                        <div class="flex-1 min-w-0">
                                            <h4
                                                class="font-extrabold text-base-content text-sm mb-1.5 truncate"
                                            >
                                                {site.nama}
                                            </h4>
                                            <p
                                                class="text-xs font-medium text-base-content/60 leading-relaxed line-clamp-2"
                                            >
                                                {site.alamat_lengkap}
                                            </p>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </section>
                {/if}

                <div class="flex justify-center pt-8 border-t border-base-200">
                    <a
                        href="/situs"
                        class="btn bg-base-100 hover:bg-base-200 border-base-200 shadow-sm rounded-xl h-14 px-8 gap-3 group"
                    >
                        <svg
                            class="w-5 h-5 text-base-content/50 transform group-hover:-translate-x-1 transition-transform"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                            ><path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M10 19l-7-7m0 0l7-7m-7 7h18"
                            /></svg
                        >
                        <span class="font-bold text-base-content/80"
                            >Kembali ke Direktori</span
                        >
                    </a>
                </div>
            </main>
        </div>
    {:else}
        <div
            class="min-h-[70vh] flex items-center justify-center relative z-10"
        >
            <div
                class="text-center max-w-md mx-auto px-4 bg-base-100/50 backdrop-blur-xl p-10 rounded-3xl border border-base-200 shadow-2xl"
            >
                <div
                    class="w-24 h-24 rounded-3xl bg-base-200/50 flex items-center justify-center mx-auto mb-6 shadow-inner"
                >
                    <svg
                        class="w-12 h-12 text-base-content/30"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                        />
                    </svg>
                </div>
                <h1
                    class="text-2xl font-extrabold tracking-tight text-base-content mb-3"
                >
                    Situs Tidak Ditemukan
                </h1>
                <p
                    class="text-sm font-medium text-base-content/60 mb-8 leading-relaxed"
                >
                    Maaf, data situs keagamaan yang Anda cari tidak tersedia,
                    URL salah, atau telah dihapus dari sistem.
                </p>
                <a
                    href="/situs"
                    class="btn btn-primary rounded-xl px-8 shadow-lg shadow-primary/20 hover:-translate-y-0.5 transition-transform font-bold"
                >
                    Cari Situs Lainnya
                </a>
            </div>
        </div>
    {/if}
</div>

<style>
    :global(.leaflet-container) {
        z-index: 0;
        font-family: inherit;
    }

    :global(.leaflet-control-attribution) {
        font-size: 10px;
        opacity: 0.8;
    }

    .line-clamp-2 {
        display: -webkit-box;
        line-clamp: 2;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
    }
</style>
