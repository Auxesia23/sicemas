<script lang="ts">
    import { onDestroy } from "svelte";
    import { goto, invalidate } from "$app/navigation";
    import { page } from "$app/state";
    import { apiFetch } from "$lib/api";
    import { auth } from "$lib/states/auth.svelte";

    import Toast from "$lib/components/ui/Toast.svelte";
    import ViewDetailMasjid from "$lib/components/views/ViewDetailMasjid.svelte";
    import ViewDetailPesantren from "$lib/components/views/ViewDetailPesantren.svelte";
    import ViewDetailMushola from "$lib/components/views/ViewDetailMushola.svelte";
    import ViewDetailMT from "$lib/components/views/ViewDetailMT.svelte";

    import type { PageData } from "./$types";

    let { data }: { data: PageData } = $props();

    let situs = $derived(data.situs);

    let verifyStatus = $derived(data.situs.status_verifikasi || "");
    let verifySitusId = $derived(data.situs.situs_id || "");
    let isVerifying = $state(false);

    let showToast = $state(false);
    let toastMessage = $state("");
    let toastType = $state<"success" | "error">("success");

    let map: any;
    let marker: any;

    onDestroy(() => {
        cleanupMap();
    });

    function mapAction(node: HTMLElement) {
        let L: any;

        (async () => {
            L = await import("leaflet");
            await import("leaflet/dist/leaflet.css");

            delete (L.Icon.Default.prototype as any)._getIconUrl;
            L.Icon.Default.mergeOptions({
                iconRetinaUrl:
                    "https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon-2x.png",
                iconUrl:
                    "https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png",
                shadowUrl:
                    "https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png",
            });

            const lat = Number(situs.latitude);
            const lng = Number(situs.longitude);

            map = L.map(node, {
                scrollWheelZoom: false,
                zoomControl: true,
                attributionControl: false,
            }).setView([lat, lng], 15);

            L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
                attribution: "© OpenStreetMap contributors",
            }).addTo(map);

            marker = L.marker([lat, lng]).addTo(map);

            setTimeout(() => {
                if (map) map.invalidateSize();
            }, 150);
        })();

        return {
            destroy() {
                cleanupMap();
            },
        };
    }

    function cleanupMap() {
        if (map) {
            map.remove();
            map = null;
            marker = null;
        }
    }

    function formatDate(dateString: string | null | undefined) {
        if (!dateString) return "-";
        const date = new Date(dateString);
        return date.toLocaleDateString("id-ID", {
            day: "numeric",
            month: "long",
            year: "numeric",
        });
    }

    async function submitVerifikasi(e: Event) {
        e.preventDefault();
        if (isVerifying) return;

        isVerifying = true;

        try {
            const res = await apiFetch(
                "/api/situs/" + page.params.id + "/verify",
                {
                    method: "PATCH",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({
                        status_verifikasi: verifyStatus,
                        situs_id: verifySitusId,
                    }),
                },
            );

            if (res.ok) {
                toastMessage = "Verifikasi berhasil disimpan";
                toastType = "success";
                showToast = true;

                invalidate("situs:detail");
                setTimeout(() => {
                    showToast = false;
                }, 2000);
            } else {
                const errText = await res.text();
                throw new Error(errText || "Gagal menyimpan verifikasi");
            }
        } catch (err: any) {
            console.error("Verification error:", err);
            toastMessage =
                err.message ||
                "Terjadi kesalahan sistem saat menyimpan verifikasi";
            toastType = "error";
            showToast = true;

            setTimeout(() => {
                showToast = false;
            }, 4000);
        } finally {
            isVerifying = false;
        }
    }

    function goBack() {
        goto("/dashboard/situs");
    }

    function getStatusBadgeClass(status: string) {
        if (!status)
            return "bg-base-200/50 text-base-content/50 border-base-200";
        switch (status.toLowerCase()) {
            case "terverifikasi":
                return "bg-success/10 text-success border-success/20";
            case "ditolak":
                return "bg-error/10 text-error border-error/20";
            case "menunggu":
                return "bg-warning/10 text-warning border-warning/20";
            default:
                return "bg-base-200/50 text-base-content/50 border-base-200";
        }
    }
</script>

<svelte:head>
    <title>{situs.nama || "Detail Situs"} | SiCemas</title>
    <meta name="description" content="Sistem Informasi Keagamaan KUA Ciemas" />
</svelte:head>

<div class="min-h-screen bg-base-200/50 p-4 sm:p-8">
    <div class="mx-auto max-w-4xl pb-10">
        <Toast
            show={showToast}
            message={toastMessage}
            type={toastType}
            onclose={() => (showToast = false)}
        />

        <div
            class="mb-8 flex flex-col gap-4 md:flex-row md:items-center md:justify-between"
        >
            <div class="flex items-start sm:items-center gap-4">
                <button
                    class="btn btn-circle btn-ghost bg-base-200/50 hover:bg-base-300 transition-colors mt-1 sm:mt-0 shadow-sm"
                    onclick={goBack}
                    title="Kembali"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="h-6 w-6"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M10 19l-7-7m0 0l7-7m-7 7h18"
                        />
                    </svg>
                </button>

                <div>
                    <div class="flex items-center gap-3 mb-2">
                        <div
                            class="w-1.5 h-8 rounded-full bg-linear-to-b from-primary to-secondary hidden sm:block"
                        ></div>
                        <h1
                            class="text-2xl sm:text-3xl font-extrabold tracking-tight text-base-content leading-tight"
                        >
                            {situs.nama}
                        </h1>
                    </div>
                    <div class="flex flex-wrap items-center gap-2 sm:ml-5">
                        <span
                            class="badge font-medium px-3 py-3 rounded-lg bg-primary/10 text-primary border-primary/20"
                        >
                            {situs.jenis_tipologi}
                        </span>
                        <span
                            class="badge font-medium px-3 py-3 rounded-lg bg-base-200/50 text-base-content/70 border-base-200"
                        >
                            {situs.jenis_situs}
                        </span>
                        {#if situs.status_verifikasi}
                            <div
                                class={`badge font-medium px-3 py-3 rounded-lg ${getStatusBadgeClass(situs.status_verifikasi)}`}
                            >
                                <div
                                    class={`w-1.5 h-1.5 rounded-full mr-1.5 ${
                                        situs.status_verifikasi.toLowerCase() ===
                                        "terverifikasi"
                                            ? "bg-success"
                                            : situs.status_verifikasi.toLowerCase() ===
                                                "menunggu"
                                              ? "bg-warning"
                                              : "bg-error"
                                    }`}
                                ></div>
                                {situs.status_verifikasi
                                    .charAt(0)
                                    .toUpperCase() +
                                    situs.status_verifikasi.slice(1)}
                            </div>
                        {/if}
                    </div>
                </div>
            </div>
            <div
                class="text-sm font-medium text-base-content/50 bg-base-200/30 px-4 py-2.5 rounded-xl border border-base-200 self-start md:self-auto shadow-sm"
            >
                Diperbarui: <span class="text-base-content/80"
                    >{formatDate(situs.updated_at)}</span
                >
            </div>
        </div>

        <div class="flex flex-col gap-6">
            <div
                class="card rounded-2xl border border-base-200 bg-base-100 shadow-xl shadow-base-200/40 overflow-hidden"
            >
                <div class="card-body p-6 sm:p-8">
                    <div
                        class="flex items-center gap-3 mb-6 border-b border-base-200 pb-4"
                    >
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
                                    d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                                /></svg
                            >
                        </div>
                        <h2
                            class="card-title text-xl font-extrabold tracking-tight text-base-content m-0"
                        >
                            Informasi Umum
                        </h2>
                    </div>

                    <div
                        class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3"
                    >
                        <div class="sm:col-span-2 lg:col-span-3">
                            <p
                                class="mb-2 text-xs font-bold uppercase tracking-wider text-base-content/50"
                            >
                                Nomor Telepon Pengurus
                            </p>
                            {#if situs.nomor_telpon_pengurus && situs.nomor_telpon_pengurus.length > 0}
                                <div class="flex flex-wrap gap-2">
                                    {#each situs.nomor_telpon_pengurus as phone}
                                        <span
                                            class="badge font-medium px-4 py-3 rounded-xl bg-primary/10 text-primary border-primary/20"
                                            >{phone}</span
                                        >
                                    {/each}
                                </div>
                            {:else}
                                <p class="font-medium text-base-content/70">
                                    -
                                </p>
                            {/if}
                        </div>
                        <div
                            class="rounded-xl bg-base-200/30 p-4 border border-base-200/50"
                        >
                            <p
                                class="mb-1 text-xs font-bold uppercase tracking-wider text-base-content/50"
                            >
                                Nomor Telepon Instansi
                            </p>
                            <p class="font-medium text-base-content">
                                {situs.nomor_telepon || "-"}
                            </p>
                        </div>
                        <div
                            class="rounded-xl bg-base-200/30 p-4 border border-base-200/50"
                        >
                            <p
                                class="mb-1 text-xs font-bold uppercase tracking-wider text-base-content/50"
                            >
                                Email
                            </p>
                            <p class="font-medium text-base-content">
                                {situs.email || "-"}
                            </p>
                        </div>
                        <div
                            class="rounded-xl bg-base-200/30 p-4 border border-base-200/50"
                        >
                            <p
                                class="mb-1 text-xs font-bold uppercase tracking-wider text-base-content/50"
                            >
                                Website
                            </p>
                            <p class="font-medium">
                                {#if situs.website}
                                    <a
                                        href={situs.website}
                                        target="_blank"
                                        class="text-primary hover:text-primary-focus hover:underline font-bold transition-colors truncate block"
                                        rel="noopener noreferrer"
                                        >{situs.website}</a
                                    >
                                {:else}
                                    <span class="text-base-content/70">-</span>
                                {/if}
                            </p>
                        </div>
                        <div
                            class="rounded-xl bg-base-200/30 p-4 border border-base-200/50"
                        >
                            <p
                                class="mb-1 text-xs font-bold uppercase tracking-wider text-base-content/50"
                            >
                                Tahun Berdiri
                            </p>
                            <p class="font-medium text-base-content">
                                {situs.tahun_berdiri || "-"}
                            </p>
                        </div>
                        <div
                            class="rounded-xl bg-base-200/30 p-4 border border-base-200/50 sm:col-span-2"
                        >
                            <p
                                class="mb-1 text-xs font-bold uppercase tracking-wider text-base-content/50"
                            >
                                Nomor Badan Hukum
                            </p>
                            <p class="font-medium text-base-content">
                                {situs.nomor_badan_hukum || "-"}
                            </p>
                        </div>
                    </div>
                </div>
            </div>

            <div
                class="card rounded-2xl border border-base-200 bg-base-100 shadow-xl shadow-base-200/40 overflow-hidden"
            >
                <div class="card-body p-6 sm:p-8">
                    <div
                        class="flex items-center gap-3 mb-6 border-b border-base-200 pb-4"
                    >
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
                                    d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"
                                /><path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"
                                /></svg
                            >
                        </div>
                        <h2
                            class="card-title text-xl font-extrabold tracking-tight text-base-content m-0"
                        >
                            Lokasi & Peta
                        </h2>
                    </div>

                    <div
                        class="mb-6 grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3"
                    >
                        <div
                            class="sm:col-span-2 lg:col-span-3 rounded-xl bg-base-200/30 p-4 border border-base-200/50"
                        >
                            <p
                                class="mb-1 text-xs font-bold uppercase tracking-wider text-base-content/50"
                            >
                                Alamat Lengkap
                            </p>
                            <p
                                class="font-medium text-base-content leading-relaxed"
                            >
                                {situs.alamat_lengkap || "-"}
                            </p>
                        </div>
                        <div
                            class="rounded-xl bg-base-200/30 p-4 border border-base-200/50"
                        >
                            <p
                                class="mb-1 text-xs font-bold uppercase tracking-wider text-base-content/50"
                            >
                                Desa/Kelurahan
                            </p>
                            <p class="font-medium text-base-content">
                                {situs.desa || "-"}
                            </p>
                        </div>
                        <div
                            class="rounded-xl bg-base-200/30 p-4 border border-base-200/50"
                        >
                            <p
                                class="mb-1 text-xs font-bold uppercase tracking-wider text-base-content/50"
                            >
                                Kecamatan
                            </p>
                            <p class="font-medium text-base-content">
                                {situs.kecamatan || "-"}
                            </p>
                        </div>
                        <div
                            class="rounded-xl bg-base-200/30 p-4 border border-base-200/50"
                        >
                            <p
                                class="mb-1 text-xs font-bold uppercase tracking-wider text-base-content/50"
                            >
                                Kabupaten/Kota
                            </p>
                            <p class="font-medium text-base-content">
                                {situs.kabupaten_kota || "-"}
                            </p>
                        </div>
                        <div
                            class="rounded-xl bg-base-200/30 p-4 border border-base-200/50"
                        >
                            <p
                                class="mb-1 text-xs font-bold uppercase tracking-wider text-base-content/50"
                            >
                                Provinsi
                            </p>
                            <p class="font-medium text-base-content">
                                {situs.provinsi || "-"}
                            </p>
                        </div>
                        <div
                            class="sm:col-span-2 rounded-xl bg-base-200/30 p-4 border border-base-200/50"
                        >
                            <p
                                class="mb-2 text-xs font-bold uppercase tracking-wider text-base-content/50"
                            >
                                Koordinat
                            </p>
                            <div
                                class="inline-flex items-center gap-2 bg-base-100 px-3 py-1.5 rounded-lg border border-base-200 shadow-sm"
                            >
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    class="h-4 w-4 text-primary"
                                    fill="none"
                                    viewBox="0 0 24 24"
                                    stroke="currentColor"
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
                                <p
                                    class="font-mono text-sm font-bold text-base-content"
                                >
                                    {situs.latitude?.toFixed(6) ?? "-"}, {situs.longitude?.toFixed(
                                        6,
                                    ) ?? "-"}
                                </p>
                            </div>
                        </div>
                    </div>

                    {#if situs.latitude && situs.longitude}
                        <div
                            use:mapAction
                            class="h-64 w-full overflow-hidden rounded-2xl border border-base-200 shadow-inner sm:h-80 z-0 relative"
                        ></div>
                    {:else}
                        <div
                            class="flex h-48 flex-col items-center justify-center rounded-2xl border-2 border-dashed border-base-300 bg-base-200/30"
                        >
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                class="h-8 w-8 text-base-content/30 mb-2"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                                ><path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M3 21v-4m0 0V5a2 2 0 012-2h6.5l1 1H21l-3 6 3 6h-8.5l-1-1H5a2 2 0 00-2 2zm9-13.5V9"
                                /></svg
                            >
                            <p class="text-sm font-bold text-base-content/40">
                                Koordinat peta tidak tersedia
                            </p>
                        </div>
                    {/if}
                </div>
            </div>

            <div
                class="card rounded-2xl border border-base-200 bg-base-100 shadow-xl shadow-base-200/40 overflow-hidden"
            >
                <div class="card-body p-6 sm:p-8">
                    <div
                        class="flex items-center gap-3 mb-6 border-b border-base-200 pb-4"
                    >
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
                                    d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                                /></svg
                            >
                        </div>
                        <h2
                            class="card-title text-xl font-extrabold tracking-tight text-base-content m-0"
                        >
                            Kapasitas & Legalitas
                        </h2>
                    </div>

                    <div
                        class="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3"
                    >
                        <div
                            class="rounded-xl bg-base-200/30 p-4 border border-base-200/50"
                        >
                            <p
                                class="mb-1 text-xs font-bold uppercase tracking-wider text-base-content/50"
                            >
                                Luas Tanah
                            </p>
                            <p class="font-medium text-base-content">
                                {situs.luas_tanah
                                    ? `${situs.luas_tanah} m²`
                                    : "-"}
                            </p>
                        </div>
                        <div
                            class="rounded-xl bg-base-200/30 p-4 border border-base-200/50"
                        >
                            <p
                                class="mb-1 text-xs font-bold uppercase tracking-wider text-base-content/50"
                            >
                                Luas Bangunan
                            </p>
                            <p class="font-medium text-base-content">
                                {situs.luas_bangunan
                                    ? `${situs.luas_bangunan} m²`
                                    : "-"}
                            </p>
                        </div>
                        <div
                            class="rounded-xl bg-base-200/30 p-4 border border-base-200/50"
                        >
                            <p
                                class="mb-1 text-xs font-bold uppercase tracking-wider text-base-content/50"
                            >
                                Status Tanah
                            </p>
                            <p class="font-medium text-base-content">
                                {situs.status_tanah || "-"}
                            </p>
                        </div>
                        <div
                            class="rounded-xl bg-base-200/30 p-4 border border-base-200/50"
                        >
                            <p
                                class="mb-1 text-xs font-bold uppercase tracking-wider text-base-content/50"
                            >
                                Daya Tampung Max
                            </p>
                            <p class="font-medium text-base-content">
                                {situs.daya_tampung_max
                                    ? `${situs.daya_tampung_max} orang`
                                    : "-"}
                            </p>
                        </div>
                        <div
                            class="rounded-xl bg-base-200/30 p-4 border border-base-200/50"
                        >
                            <p
                                class="mb-1 text-xs font-bold uppercase tracking-wider text-base-content/50"
                            >
                                Nomor AIW
                            </p>
                            <p class="font-medium text-base-content">
                                {situs.nomor_aiw || "-"}
                            </p>
                        </div>
                        <div
                            class="rounded-xl bg-base-200/30 p-4 border border-base-200/50"
                        >
                            <p
                                class="mb-1 text-xs font-bold uppercase tracking-wider text-base-content/50"
                            >
                                No. Sertifikat Wakaf
                            </p>
                            <p class="font-medium text-base-content">
                                {situs.nomor_sertifikat_wakaf || "-"}
                            </p>
                        </div>
                    </div>
                </div>
            </div>

            {#if situs.foto && situs.foto.length > 0}
                <div
                    class="card rounded-2xl border border-base-200 bg-base-100 shadow-xl shadow-base-200/40 overflow-hidden"
                >
                    <div class="card-body p-6 sm:p-8">
                        <div
                            class="flex items-center gap-3 mb-6 border-b border-base-200 pb-4"
                        >
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
                                        d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
                                    /></svg
                                >
                            </div>
                            <h2
                                class="card-title text-xl font-extrabold tracking-tight text-base-content m-0"
                            >
                                Foto Situs
                            </h2>
                        </div>
                        <div
                            class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4"
                        >
                            {#each situs.foto as foto (foto.id)}
                                <div
                                    class="group overflow-hidden rounded-xl border border-base-200 shadow-sm relative"
                                >
                                    <img
                                        src={foto.image_url}
                                        alt="Foto situs"
                                        class="h-48 w-full object-cover transition-transform duration-500 group-hover:scale-110"
                                    />
                                    <div
                                        class="absolute inset-0 bg-black/0 transition-colors duration-300 group-hover:bg-black/10"
                                    ></div>
                                </div>
                            {/each}
                        </div>
                    </div>
                </div>
            {/if}

            {#if situs.jenis_situs === "Masjid" || situs.jenis_situs
                    .toLowerCase()
                    .includes("masjid")}
                <ViewDetailMasjid detail={situs.detail} />
            {:else if situs.jenis_situs === "Ponpes" || situs.jenis_situs
                    .toLowerCase()
                    .includes("ponpes")}
                <ViewDetailPesantren detail={situs.detail} />
            {:else if situs.jenis_situs === "Mushola" || situs.jenis_situs
                    .toLowerCase()
                    .includes("mushola")}
                <ViewDetailMushola detail={situs.detail} />
            {:else if situs.jenis_situs === "MT" || situs.jenis_situs
                    .toLowerCase()
                    .includes("majelis")}
                <ViewDetailMT detail={situs.detail} />
            {/if}

            {#if auth.hasAllPermissions(["situs:verify"])}
                <div
                    class="card rounded-2xl border border-primary/20 bg-primary/5 shadow-xl shadow-primary/10 mt-2"
                >
                    <div class="card-body p-6 sm:p-8">
                        <div
                            class="flex items-center gap-3 border-b border-primary/10 pb-4 mb-6"
                        >
                            <div
                                class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-linear-to-br from-primary/30 to-secondary/20 text-primary shadow-sm"
                            >
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    class="h-5 w-5 font-bold"
                                    fill="none"
                                    viewBox="0 0 24 24"
                                    stroke="currentColor"
                                    ><path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                                    /></svg
                                >
                            </div>
                            <div>
                                <h2
                                    class="card-title text-xl font-extrabold text-primary m-0 tracking-tight"
                                >
                                    Panel Verifikasi Situs
                                </h2>
                                <p
                                    class="text-xs font-medium text-primary/70 mt-0.5"
                                >
                                    Tinjau dan ubah status verifikasi data situs
                                    ini.
                                </p>
                            </div>
                        </div>

                        <form onsubmit={submitVerifikasi} class="space-y-6">
                            <div class="grid grid-cols-1 gap-6 sm:grid-cols-2">
                                <div class="form-control">
                                    <label
                                        class="label pb-2 pt-0"
                                        for="verifyStatus"
                                    >
                                        <span
                                            class="label-text font-bold text-base-content/80"
                                            >Status Verifikasi <span
                                                class="text-error">*</span
                                            ></span
                                        >
                                    </label>
                                    <select
                                        id="verifyStatus"
                                        class="select select-bordered w-full rounded-xl bg-base-100 min-h-12 focus:border-primary transition-colors shadow-sm"
                                        required
                                        bind:value={verifyStatus}
                                    >
                                        <option value="" disabled
                                            >Pilih status...</option
                                        >
                                        <option value="terverifikasi"
                                            >Terverifikasi</option
                                        >
                                        <option value="ditolak">Ditolak</option>
                                    </select>
                                </div>

                                <div class="form-control">
                                    <label
                                        class="label pb-2 pt-0"
                                        for="verifySitusId"
                                    >
                                        <span
                                            class="label-text font-bold text-base-content/80"
                                        >
                                            Nomor Statistik Kemenag <span
                                                class="text-xs font-normal opacity-50"
                                                >(Opsional)</span
                                            >
                                        </span>
                                    </label>
                                    <input
                                        id="verifySitusId"
                                        type="text"
                                        class="input input-bordered w-full min-h-12 rounded-xl bg-base-100 focus:border-primary transition-colors shadow-sm disabled:bg-base-200/50"
                                        placeholder="Contoh: 51.2.12.345678"
                                        bind:value={verifySitusId}
                                        disabled={verifyStatus === "ditolak" ||
                                            verifyStatus === "menunggu"}
                                    />
                                    {#if verifyStatus === "ditolak" || verifyStatus === "menunggu"}
                                        <span
                                            class="mt-1.5 text-xs font-medium text-warning"
                                            >Tidak dapat diisi saat status
                                            ditolak/menunggu</span
                                        >
                                    {/if}
                                </div>
                            </div>

                            <div
                                class="flex justify-end pt-4 border-t border-primary/10"
                            >
                                <button
                                    type="submit"
                                    class="btn btn-primary rounded-xl px-8 min-h-12 shadow-lg shadow-primary/20 transition-transform hover:-translate-y-0.5 w-full sm:w-auto font-bold"
                                    disabled={isVerifying || !verifyStatus}
                                >
                                    {#if isVerifying}
                                        <span
                                            class="loading loading-sm loading-spinner mr-2"
                                        ></span> Menyimpan...
                                    {:else}
                                        Simpan Verifikasi
                                    {/if}
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            {/if}
        </div>
    </div>
</div>

<style>
    :global(.leaflet-container) {
        z-index: 0;
    }
</style>
