<script>
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { navigating } from '$app/stores';
	import { page } from '$app/stores';

	// Receive SSR data via props
	let { data } = $props();

	// Derived state - always in sync with server response
	let sites = $derived(data.sites || []);

	// Client-side state
	let searchTerm = $state('');
	let isLocating = $state(false);
	let userLocation = $state(null);

	// Derived loading state - combines navigation + GPS locating
	let isLoading = $derived($navigating !== null || isLocating);

	// Derived filtered sites
	let filteredSites = $derived(
		sites.filter((site) => {
			const searchLower = searchTerm.toLowerCase();
			return (
				site.nama.toLowerCase().includes(searchLower) ||
				site.jenis.toLowerCase().includes(searchLower) ||
				site.desa.toLowerCase().includes(searchLower)
			);
		})
	);

	// Format distance helper
	function formatDistance(meters) {
		if (meters === null || meters === undefined) return null;
		if (meters < 1000) {
			return `${Math.round(meters)}m`;
		}
		return `${(meters / 1000).toFixed(1)}km`;
	}

	// GPS Location handler - called once on mount
	onMount(async () => {
		// Check if we already have location in URL - prevent infinite loop
		const existingLat = $page.url.searchParams.get('lat');
		const existingLng = $page.url.searchParams.get('lng');

		if (existingLat && existingLng) {
			// Already have location, no need to request GPS again
			userLocation = { lat: parseFloat(existingLat), lng: parseFloat(existingLng) };
			return;
		}

		// No location in URL, request GPS
		isLocating = true;

		try {
			// Request geolocation with timeout
			const position = await new Promise((resolve, reject) => {
				const timeoutId = setTimeout(() => {
					reject(new Error('GPS timeout'));
				}, 10000);

				navigator.geolocation.getCurrentPosition(
					(position) => {
						clearTimeout(timeoutId);
						resolve(position);
					},
					(error) => {
						clearTimeout(timeoutId);
						reject(error);
					},
					{
						enableHighAccuracy: true,
						timeout: 10000,
						maximumAge: 0
					}
				);
			});

			// Got GPS coordinates
			const lat = position.coords.latitude;
			const lng = position.coords.longitude;

			// Update local state
			userLocation = { lat, lng };

			// Navigate to update URL and trigger SSR re-fetch
			// Use replaceState to avoid breaking back button
			await goto(`?lat=${lat}&lng=${lng}`, {
				replaceState: true,
				noScroll: true
			});
		} catch (geoError) {
			// GPS denied, error, or timeout - that's okay, show unsorted list
			console.warn('GPS not available:', geoError.message);
			userLocation = null;
		} finally {
			// Always reset loading state
			isLocating = false;
		}
	});
</script>

<svelte:head>
	<title>Direktori Situs Keagamaan - KUA Ciemas</title>
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

	<!-- Progress Bar (shows during navigation/GPS) -->
	{#if isLoading}
		<div class="fixed top-0 left-0 z-50 h-1 w-full bg-primary/20">
			<div class="h-full w-full animate-pulse bg-primary"></div>
		</div>
	{/if}

	<!-- Navigation -->
	<nav class="navbar bg-base-100/80 px-4 backdrop-blur-sm sm:px-6 lg:px-8">
		<div class="flex-1">
			<a href="/" class="btn text-xl font-bold text-primary btn-ghost">SICEMAS</a>
		</div>
		<div class="flex-none gap-2">
			{#if userLocation}
				<span class="badge badge-sm badge-success">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="mr-1 h-3 w-3"
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
					Lokasi Aktif
				</span>
			{/if}
			<a href="/login" class="btn btn-primary">Masuk</a>
		</div>
	</nav>

	<!-- Main Content -->
	<main class="relative container mx-auto px-4 py-12 sm:px-6 lg:px-8">
		<!-- Hero Section -->
		<section class="mb-20 text-center">
			<span class="mb-6 badge badge-lg badge-primary">KUA Kecamatan Ciemas</span>
			<h1 class="mb-6 text-3xl font-bold text-gray-800 md:text-5xl">
				Direktori <span
					class="bg-linear-to-r from-primary to-secondary bg-clip-text text-transparent"
					>Situs Keagamaan</span
				>
			</h1>
			<p class="mx-auto mb-12 max-w-2xl text-lg text-gray-600">
				Temukan dan jelajahi situs-situs keagamaan di wilayah Kecamatan Ciemas, Kabupaten Sukabumi
			</p>

			<!-- Search Bar -->
			<div class="mx-auto max-w-xl">
				<div class="form-control">
					<div class="input-group relative">
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="absolute top-1/2 left-3 h-5 w-5 -translate-y-1/2 text-gray-400"
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
							class="input-bordered input w-full pl-12 shadow-lg"
							placeholder="Cari nama situs, jenis, atau lokasi..."
							bind:value={searchTerm}
						/>
					</div>
				</div>
			</div>
		</section>

		<!-- Loading State -->
		{#if isLoading}
			<div class="flex items-center justify-center py-20">
				<div class="flex flex-col items-center gap-4">
					<span class="loading loading-lg loading-spinner text-primary"></span>
					<p class="text-gray-600">
						{isLocating ? 'Mengambil lokasi Anda...' : 'Memuat data situs...'}
					</p>
				</div>
			</div>
		{:else if sites.length === 0}
			<!-- Empty State -->
			<div class="flex items-center justify-center py-20">
				<div class="flex flex-col items-center gap-4 text-center">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="h-16 w-16 text-gray-400"
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
					<h3 class="text-xl font-bold text-gray-800">Belum Ada Data</h3>
					<p class="text-gray-600">Belum ada situs keagamaan yang terdaftar</p>
				</div>
			</div>
		{:else if filteredSites.length === 0}
			<!-- Search Empty State -->
			<div class="flex items-center justify-center py-20">
				<div class="flex flex-col items-center gap-4 text-center">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="h-16 w-16 text-gray-400"
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
					<h3 class="text-xl font-bold text-gray-800">Tidak Ditemukan</h3>
					<p class="text-gray-600">
						{searchTerm
							? `Tidak ada situs yang cocok dengan "${searchTerm}"`
							: 'Tidak ada hasil pencarian'}
					</p>
					{#if searchTerm}
						<button class="btn btn-outline" onclick={() => (searchTerm = '')}>
							Hapus Pencarian
						</button>
					{/if}
				</div>
			</div>
		{:else}
			<!-- Sites Grid -->
			<div class="grid grid-cols-1 gap-8 md:grid-cols-2 lg:grid-cols-3">
				{#each filteredSites as site (site.id)}
					<div
						class="group card border border-base-200 bg-base-100 shadow-lg transition-all duration-300 hover:-translate-y-2 hover:shadow-2xl"
					>
						<!-- Image Section with Fixed Height -->
						<figure class="relative h-56 w-full overflow-hidden">
							{#if site.foto_utama}
								<img
									src={site.foto_utama}
									alt={site.nama}
									class="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
								/>
							{:else}
								<div
									class="flex h-full w-full items-center justify-center bg-linear-to-br from-primary/20 to-secondary/20"
								>
									<svg
										xmlns="http://www.w3.org/2000/svg"
										class="h-20 w-20 text-gray-400"
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
								</div>
							{/if}

							<!-- Overlay Badges -->
							<div class="absolute top-2 left-2 flex flex-col gap-2">
								<span class="badge text-xs badge-primary">
									{site.jenis}
								</span>
								{#if site.jarak_meter}
									<span class="badge text-xs badge-secondary">
										<svg
											xmlns="http://www.w3.org/2000/svg"
											class="mr-1 h-3 w-3"
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
						</figure>

						<!-- Card Body -->
						<div class="card-body flex flex-col justify-between p-4">
							<div>
								<h3 class="card-title text-lg font-bold text-gray-800">{site.nama}</h3>
								<div class="mt-2 flex items-start gap-2 text-sm text-gray-600">
									<svg
										xmlns="http://www.w3.org/2000/svg"
										class="mt-0.5 h-4 w-4 shrink-0"
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
									<span>{site.desa}, Kecamatan Ciemas</span>
								</div>
								{#if site.jarak_meter}
									<div class="mt-1 text-xs text-gray-500">
										Jarak: {formatDistance(site.jarak_meter)} dari lokasi Anda
									</div>
								{/if}
							</div>

							<!-- Action Button -->
							<div class="mt-4">
								<a href="/situs/{site.id}" class="btn w-full btn-outline btn-sm btn-primary">
									Lihat Detail
								</a>
							</div>
						</div>
					</div>
				{/each}
			</div>
		{/if}

		<!-- Stats Section -->
		{#if !isLoading && sites.length > 0}
			<section class="mt-24">
				<div class="rounded-2xl bg-base-100 p-8 shadow-xl">
					<div class="grid grid-cols-2 gap-8 text-center md:grid-cols-4">
						<div>
							<div class="text-3xl font-bold text-primary">{sites.length}</div>
							<div class="text-sm text-gray-500">Total Situs</div>
						</div>
						<div>
							<div class="text-3xl font-bold text-secondary">
								{new Set(sites.map((s) => s.jenis)).size}
							</div>
							<div class="text-sm text-gray-500">Jenis Situs</div>
						</div>
						<div>
							<div class="text-3xl font-bold text-accent">
								{sites.filter((s) => s.foto_utama).length}
							</div>
							<div class="text-sm text-gray-500">Dengan Foto</div>
						</div>
						<div>
							<div class="text-3xl font-bold text-primary">100%</div>
							<div class="text-sm text-gray-500">Terverifikasi</div>
						</div>
					</div>
				</div>
			</section>
		{/if}
	</main>

	<!-- Footer -->
	<footer class="footer-center mt-24 footer rounded bg-base-200 p-10 text-base-content">
		<div>
			<div class="mb-4 flex items-center gap-2">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="h-8 w-8 text-primary"
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
