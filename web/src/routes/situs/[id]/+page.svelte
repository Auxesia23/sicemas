<script>
	import { onMount } from 'svelte';
	import { page } from '$app/stores';

	// Receive SSR data via props
	let { data } = $props();

	// Derived state - always in sync with server
	let site = $derived(data.site);

	// Map state
	let mapContainer = $state();
	let map;
	let marker;

	// Format distance helper
	function formatDistance(meters) {
		if (meters === null || meters === undefined) return '-';
		if (meters < 1000) {
			return `${Math.round(meters)}m`;
		}
		return `${(meters / 1000).toFixed(1)}km`;
	}

	// Format key to readable text
	function formatKey(key) {
		return key
			.split('_')
			.map((word) => word.charAt(0).toUpperCase() + word.slice(1))
			.join(' ');
	}

	// Initialize map after component mounts
	onMount(async () => {
		if (!site.latitude || !site.longitude) return;

		// Dynamic import Leaflet (client-side only)
		const L = await import('leaflet');
		await import('leaflet/dist/leaflet.css');

		delete L.Icon.Default.prototype._getIconUrl;
		L.Icon.Default.mergeOptions({
			iconRetinaUrl:
				'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-icon-2x.png',
			iconUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-icon.png',
			shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-shadow.png'
		});

		// Initialize map
		map = L.map(mapContainer, {
			dragging: true,
			touchZoom: true,
			scrollWheelZoom: false,
			doubleClickZoom: false,
			zoomControl: true,
			attributionControl: true
		}).setView([site.latitude, site.longitude], 15);

		// Add tile layer
		L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
			attribution: '© OpenStreetMap contributors'
		}).addTo(map);

		// Add marker
		marker = L.marker([site.latitude, site.longitude]).addTo(map);

		// Fix map rendering
		setTimeout(() => {
			map.invalidateSize();
		}, 200);
	});
</script>

<svelte:head>
	<title>{site?.nama || 'Detail Situs'} - KUA Ciemas</title>
	<meta name="description" content={site?.alamat_lengkap} />
</svelte:head>

<div class="min-h-screen bg-linear-to-b from-base-100 to-base-200">
	<!-- Animated Background Blobs -->
	<div class="pointer-events-none fixed inset-0 overflow-hidden">
		<div
			class="animate-blob absolute top-20 -left-4 h-72 w-72 rounded-full bg-primary opacity-10 mix-blend-multiply blur-3xl"
		></div>
		<div
			class="animate-blob animation-delay-2000 absolute top-40 -right-4 h-72 w-72 rounded-full bg-secondary opacity-10 mix-blend-multiply blur-3xl"
		></div>
		<div
			class="animate-blob animation-delay-4000 absolute bottom-20 left-1/2 h-72 w-72 rounded-full bg-accent opacity-10 mix-blend-multiply blur-3xl"
		></div>
	</div>

	<!-- Navigation -->
	<nav class="navbar bg-base-100/80 px-4 backdrop-blur-sm sm:px-6 lg:px-8">
		<div class="flex-1">
			<a href="/" class="btn text-xl font-bold text-primary btn-ghost"
				><img src="/kemenag.webp" alt="SICEMAS" class="mr-2 h-8" />SICEMAS</a
			>
		</div>
		<div class="flex-none">
			{#if data.user}
				<a href="/dashboard" class="btn btn-primary">Dashboard</a>
			{:else}
				<a href="/login" class="btn btn-primary">Masuk</a>
			{/if}
		</div>
	</nav>

	<!-- Main Content -->
	<main class="relative container mx-auto px-4 py-12 sm:px-6 lg:px-8">
		{#if site}
			<!-- Back Button -->
			<div class="mb-6">
				<a href="/situs" class="btn btn-ghost btn-sm">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="mr-2 h-4 w-4"
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
					Kembali ke Direktori
				</a>
			</div>

			<!-- Header Section -->
			<section class="mb-8">
				<div class="flex flex-col gap-4">
					<div>
						<h1 class="text-3xl font-bold text-gray-800 md:text-4xl">{site.nama}</h1>
						<div class="mt-3 flex flex-wrap items-center gap-2">
							<span class="badge badge-lg badge-primary">{site.jenis}</span>
							<span class="badge badge-lg badge-secondary">
								<svg
									xmlns="http://www.w3.org/2000/svg"
									class="mr-1 h-4 w-4"
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
								{site.desa}, Kecamatan Ciemas
							</span>
							{#if site.jarak_meter}
								<span class="badge badge-lg badge-accent">
									📍 {formatDistance(site.jarak_meter)}
								</span>
							{/if}
						</div>
					</div>
					<div class="flex items-start gap-2 text-gray-600">
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="mt-1 h-5 w-5 shrink-0"
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
						<p class="text-lg">{site.alamat_lengkap}</p>
					</div>
				</div>
			</section>

			<!-- Stats Section -->
			<section class="mb-8">
				<div class="grid grid-cols-2 gap-4 md:grid-cols-4">
					{#if site.luas_tanah}
						<div class="card border border-base-200 bg-base-100 p-4 shadow-md">
							<div class="text-sm text-gray-500">Luas Tanah</div>
							<div class="text-2xl font-bold text-primary">{site.luas_tanah} m²</div>
						</div>
					{/if}
					{#if site.luas_bangunan}
						<div class="card border border-base-200 bg-base-100 p-4 shadow-md">
							<div class="text-sm text-gray-500">Luas Bangunan</div>
							<div class="text-2xl font-bold text-secondary">{site.luas_bangunan} m²</div>
						</div>
					{/if}
					{#if site.daya_tampung}
						<div class="card border border-base-200 bg-base-100 p-4 shadow-md">
							<div class="text-sm text-gray-500">Daya Tampung</div>
							<div class="text-2xl font-bold text-accent">{site.daya_tampung} Jamaah</div>
						</div>
					{/if}
					{#if site.tahun_berdiri}
						<div class="card border border-base-200 bg-base-100 p-4 shadow-md">
							<div class="text-sm text-gray-500">Tahun Berdiri</div>
							<div class="text-2xl font-bold text-primary">{site.tahun_berdiri}</div>
						</div>
					{/if}
				</div>
			</section>

			<!-- Gallery Section -->
			{#if site.galeri && site.galeri.length > 0}
				<section class="mb-8">
					<h2 class="mb-4 text-2xl font-bold text-gray-800">Galeri Foto</h2>
					<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
						<!-- Hero Image (First Photo) -->
						<div class="relative overflow-hidden rounded-xl md:col-span-2 lg:col-span-2">
							<img src={site.galeri[0]} alt={site.nama} class="h-64 w-full object-cover md:h-96" />
						</div>
						<!-- Remaining Photos Grid -->
						{#if site.galeri.length > 1}
							<div class="grid grid-cols-2 gap-4 md:grid-cols-1 lg:grid-cols-2">
								{#each site.galeri.slice(1, 5) as foto (foto)}
									<div class="relative overflow-hidden rounded-xl">
										<img src={foto} alt="Galeri situs" class="h-32 w-full object-cover md:h-48" />
									</div>
								{/each}
							</div>
						{/if}
					</div>
				</section>
			{:else}
				<!-- Empty Gallery State -->
				<section class="mb-8">
					<div class="rounded-xl bg-base-100 p-12 text-center shadow-md">
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="mx-auto mb-4 h-16 w-16 text-gray-400"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
							/>
						</svg>
						<h3 class="text-lg font-bold text-gray-800">Belum Ada Foto</h3>
						<p class="text-gray-600">Foto untuk situs ini belum ditambahkan</p>
					</div>
				</section>
			{/if}

			<!-- Fasilitas Section (Dynamic) -->
			{#if site.fasilitas && Object.keys(site.fasilitas).length > 0}
				<section class="mb-8">
					<h2 class="mb-4 text-2xl font-bold text-gray-800">Fasilitas</h2>
					<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
						{#each Object.entries(site.fasilitas) as [key, values]}
							{#if values && values.length > 0}
								<div class="card border border-base-200 bg-base-100 p-4 shadow-md">
									<h3 class="mb-3 text-lg font-bold text-gray-800">{formatKey(key)}</h3>
									<div class="flex flex-wrap gap-2">
										{#each values as item (item)}
											<span class="badge badge-outline">{item}</span>
										{/each}
									</div>
								</div>
							{/if}
						{/each}
					</div>
				</section>
			{/if}

			<!-- Map Section -->
			{#if site.latitude && site.longitude}
				<section class="mb-8">
					<h2 class="mb-4 text-2xl font-bold text-gray-800">Lokasi</h2>
					<div class="overflow-hidden rounded-xl border border-base-200 shadow-lg">
						<div id="map-container" bind:this={mapContainer} class="h-80 w-full"></div>
					</div>
					<div class="mt-4 flex justify-center">
						<a
							href="https://www.google.com/maps/dir/?api=1&destination={site.latitude},{site.longitude}"
							target="_blank"
							rel="noopener noreferrer"
							class="btn btn-primary"
						>
							<svg
								xmlns="http://www.w3.org/2000/svg"
								class="mr-2 h-5 w-5"
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
							Buka di Google Maps
						</a>
					</div>
				</section>
			{/if}
		{:else}
			<!-- Loading State -->
			<div class="flex items-center justify-center py-20">
				<div class="flex flex-col items-center gap-4">
					<span class="loading loading-lg loading-spinner text-primary"></span>
					<p class="text-gray-600">Memuat detail situs...</p>
				</div>
			</div>
		{/if}
	</main>

	<!-- Footer -->
	<footer class="footer-center mt-24 footer rounded bg-base-200 p-10 text-base-content">
		<div>
			<div class="mb-4 flex items-center gap-2">
				<img src="/kemenag.webp" alt="SICEMAS" class="mr-2 h-8" />
				<span class="text-lg font-bold">SICEMAS</span>
			</div>
		</div>
	</footer>
</div>

<style>
	@keyframes blob {
		0% {
			transform: translate(0px, 0px) scale(1);
		}
		33% {
			transform: translate(30px, -50px) scale(1.1);
		}
		66% {
			transform: translate(-20px, 20px) scale(0.9);
		}
		100% {
			transform: translate(0px, 0px) scale(1);
		}
	}

	.animate-blob {
		animation: blob 7s infinite;
	}

	.animation-delay-2000 {
		animation-delay: 2s;
	}

	.animation-delay-4000 {
		animation-delay: 4s;
	}
</style>
