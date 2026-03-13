<script>
	import { onMount, onDestroy } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import apiService from '$lib/api';

	import ViewDetailMasjid from '$lib/components/views/ViewDetailMasjid.svelte';
	import ViewDetailPesantren from '$lib/components/views/ViewDetailPesantren.svelte';
	import ViewDetailMushola from '$lib/components/views/ViewDetailMushola.svelte';
	import ViewDetailMT from '$lib/components/views/ViewDetailMT.svelte';

	let situs = $state(null);
	let isLoading = $state(true);
	let error = $state(null);

	let map;
	let marker;

	onMount(async () => {
		try {
			isLoading = true;
			error = null;

			const response = await apiService.get(`/situs/${$page.params.id}`);

			if (!response.ok) {
				throw new Error('Gagal memuat data situs');
			}

			situs = await response.json();
		} catch (err) {
			console.error(err);
			error = err.message;
		} finally {
			isLoading = false;
		}
	});

	onDestroy(() => {
		cleanupMap();
	});

	// Action Svelte untuk Leaflet
	function mapAction(node) {
		let L;

		(async () => {
			L = await import('leaflet');
			await import('leaflet/dist/leaflet.css');

			// FIX marker icon path
			delete L.Icon.Default.prototype._getIconUrl;
			L.Icon.Default.mergeOptions({
				iconRetinaUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon-2x.png',
				iconUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png',
				shadowUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png'
			});

			const lat = Number(situs.latitude);
			const lng = Number(situs.longitude);

			map = L.map(node, {
				dragging: false,
				touchZoom: false,
				scrollWheelZoom: false,
				doubleClickZoom: false,
				boxZoom: false,
				keyboard: false,
				zoomControl: false,
				attributionControl: false
			}).setView([lat, lng], 15);

			L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
				attribution: '© OpenStreetMap contributors'
			}).addTo(map);

			marker = L.marker([lat, lng]).addTo(map);

			// Beri jeda 150ms agar Tailwind merender ukuran div dengan sempurna
			setTimeout(() => {
				if (map) map.invalidateSize();
			}, 150);
		})();

		return {
			destroy() {
				cleanupMap();
			}
		};
	}

	function cleanupMap() {
		if (map) {
			map.remove();
			map = null;
			marker = null;
		}
	}

	function formatDate(dateString) {
		if (!dateString) return '-';

		const date = new Date(dateString);

		return date.toLocaleDateString('id-ID', {
			day: 'numeric',
			month: 'long',
			year: 'numeric'
		});
	}

	function goBack() {
		goto('/dashboard/situs');
	}
</script>

<svelte:head>
	<title>{situs?.nama || 'Detail Situs'} - Dashboard</title>
</svelte:head>

<div class="mx-auto max-w-7xl p-4">
	{#if isLoading}
		<div class="flex min-h-[60vh] items-center justify-center">
			<div class="flex flex-col items-center gap-4">
				<span class="loading loading-lg loading-spinner text-primary"></span>
				<p class="text-base-content/70">Memuat data situs...</p>
			</div>
		</div>
	{:else if error}
		<div class="flex min-h-[60vh] flex-col items-center justify-center gap-4">
			<svg
				xmlns="http://www.w3.org/2000/svg"
				class="h-16 w-16 text-error"
				fill="none"
				viewBox="0 0 24 24"
				stroke="currentColor"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="2"
					d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
				/>
			</svg>
			<p class="text-error">{error}</p>
			<button class="btn btn-primary" onclick={goBack}>Kembali</button>
		</div>
	{:else if situs}
		<div class="space-y-6">
			<div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
				<div class="flex items-center gap-3">
					<button class="btn btn-square btn-ghost" onclick={goBack} title="Kembali">
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
						<h1 class="text-2xl font-bold text-primary sm:text-3xl">{situs.nama}</h1>
						<div class="mt-1 flex items-center gap-2">
							<span class="badge badge-lg badge-primary">{situs.jenis_tipologi}</span>
							<span class="badge badge-outline">{situs.jenis_situs}</span>
						</div>
					</div>
				</div>
				<div class="text-sm text-base-content/70">
					Diperbarui: {formatDate(situs.updated_at)}
				</div>
			</div>

			<div class="card border border-base-200 bg-base-100 shadow-lg">
				<div class="card-body p-4 sm:p-6">
					<h2 class="mb-4 card-title text-lg">Informasi Umum</h2>
					<div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
						<div>
							<p class="mb-1 block text-xs font-medium text-base-content/70">ID (Kemenag)</p>
							<p class="font-mono font-medium text-primary">
								{situs.situs_id || '-'}
							</p>
						</div>
						<div>
							<p class="mb-1 block text-xs font-medium text-base-content/70">Email</p>
							<p class="font-medium">{situs.email || '-'}</p>
						</div>
						<div>
							<p class="mb-1 block text-xs font-medium text-base-content/70">Nomor Telepon</p>
							<p class="font-medium">{situs.nomor_telepon || '-'}</p>
						</div>
						<div>
							<p class="mb-1 block text-xs font-medium text-base-content/70">Website</p>
							<p class="font-medium">
								{#if situs.website}
									<a
										href={situs.website}
										target="_blank"
										class="link link-primary"
										rel="noopener noreferrer"
									>
										{situs.website}
									</a>
								{:else}
									-
								{/if}
							</p>
						</div>
						<div>
							<p class="mb-1 block text-xs font-medium text-base-content/70">Tahun Berdiri</p>
							<p class="font-medium">{situs.tahun_berdiri || '-'}</p>
						</div>
						<div>
							<p class="mb-1 block text-xs font-medium text-base-content/70">Nomor Badan Hukum</p>
							<p class="font-medium">{situs.nomor_badan_hukum || '-'}</p>
						</div>
					</div>
				</div>
			</div>

			<div class="card border border-base-200 bg-base-100 shadow-lg">
				<div class="card-body p-4 sm:p-6">
					<h2 class="mb-4 card-title text-lg">Lokasi & Peta</h2>
					<div class="mb-4 grid grid-cols-1 gap-4 sm:grid-cols-2">
						<div>
							<p class="mb-1 block text-xs font-medium text-base-content/70">Alamat Lengkap</p>
							<p class="font-medium">{situs.alamat_lengkap || '-'}</p>
						</div>
						<div>
							<p class="mb-1 block text-xs font-medium text-base-content/70">Desa/Kelurahan</p>
							<p class="font-medium">{situs.desa || '-'}</p>
						</div>
						<div>
							<p class="mb-1 block text-xs font-medium text-base-content/70">Kecamatan</p>
							<p class="font-medium">{situs.kecamatan || '-'}</p>
						</div>
						<div>
							<p class="mb-1 block text-xs font-medium text-base-content/70">Kabupaten/Kota</p>
							<p class="font-medium">{situs.kabupaten_kota || '-'}</p>
						</div>
						<div>
							<p class="mb-1 block text-xs font-medium text-base-content/70">Provinsi</p>
							<p class="font-medium">{situs.provinsi || '-'}</p>
						</div>
						<div>
							<p class="mb-1 block text-xs font-medium text-base-content/70">Koordinat</p>
							<p class="font-mono text-sm font-medium">
								{situs.latitude?.toFixed(6) ?? '-'}, {situs.longitude?.toFixed(6) ?? '-'}
							</p>
						</div>
					</div>

					{#if situs.latitude && situs.longitude}
						<div
							use:mapAction
							class="h-64 w-full overflow-hidden rounded-lg border border-base-300 sm:h-80"
						></div>
					{:else}
						<div class="flex h-48 items-center justify-center rounded-lg bg-base-200">
							<p class="text-base-content/50">Koordinat tidak tersedia</p>
						</div>
					{/if}
				</div>
			</div>

			<div class="card border border-base-200 bg-base-100 shadow-lg">
				<div class="card-body p-4 sm:p-6">
					<h2 class="mb-4 card-title text-lg">Kapasitas & Legalitas</h2>
					<div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
						<div>
							<p class="mb-1 block text-xs font-medium text-base-content/70">Luas Tanah</p>
							<p class="font-medium">{situs.luas_tanah ? `${situs.luas_tanah} m²` : '-'}</p>
						</div>
						<div>
							<p class="mb-1 block text-xs font-medium text-base-content/70">Luas Bangunan</p>
							<p class="font-medium">{situs.luas_bangunan ? `${situs.luas_bangunan} m²` : '-'}</p>
						</div>
						<div>
							<p class="mb-1 block text-xs font-medium text-base-content/70">Status Tanah</p>
							<p class="font-medium">{situs.status_tanah || '-'}</p>
						</div>
						<div>
							<p class="mb-1 block text-xs font-medium text-base-content/70">Daya Tampung Max</p>
							<p class="font-medium">
								{situs.daya_tampung_max ? `${situs.daya_tampung_max} orang` : '-'}
							</p>
						</div>
						<div>
							<p class="mb-1 block text-xs font-medium text-base-content/70">Nomor AIW</p>
							<p class="font-medium">{situs.nomor_aiw || '-'}</p>
						</div>
						<div>
							<p class="mb-1 block text-xs font-medium text-base-content/70">
								No. Sertifikat Wakaf
							</p>
							<p class="font-medium">{situs.nomor_sertifikat_wakaf || '-'}</p>
						</div>
					</div>
				</div>
			</div>

			{#if situs.foto && situs.foto.length > 0}
				<div class="card border border-base-200 bg-base-100 shadow-lg">
					<div class="card-body p-4 sm:p-6">
						<h2 class="mb-4 card-title text-lg">Foto Situs</h2>
						<div class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4">
							{#each situs.foto as foto (foto.id)}
								<div class="overflow-hidden rounded-lg border border-base-300">
									<img
										src={foto.image_url}
										alt="Foto situs"
										class="h-48 w-full object-cover transition-transform duration-300 hover:scale-110"
									/>
								</div>
							{/each}
						</div>
					</div>
				</div>
			{/if}

			{#if situs.jenis_situs === 'Masjid' || situs.jenis_situs.toLowerCase().includes('masjid')}
				<ViewDetailMasjid detail={situs.detail} />
			{:else if situs.jenis_situs === 'Ponpes' || situs.jenis_situs
					.toLowerCase()
					.includes('ponpes')}
				<ViewDetailPesantren detail={situs.detail} />
			{:else if situs.jenis_situs === 'Mushola' || situs.jenis_situs
					.toLowerCase()
					.includes('mushola')}
				<ViewDetailMushola detail={situs.detail} />
			{:else if situs.jenis_situs === 'MT' || situs.jenis_situs.toLowerCase().includes('majelis')}
				<ViewDetailMT detail={situs.detail} />
			{/if}
		</div>
	{/if}
</div>

<style>
	:global(.leaflet-container) {
		z-index: 0;
	}
</style>
