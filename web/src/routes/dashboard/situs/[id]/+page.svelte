<script>
	import { onMount, onDestroy } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import apiService from '$lib/api';
	import { hasAnyPermission } from '$lib/permissions';

	// UI Components
	import Toast from '$lib/components/ui/Toast.svelte';
	import ViewDetailMasjid from '$lib/components/views/ViewDetailMasjid.svelte';
	import ViewDetailPesantren from '$lib/components/views/ViewDetailPesantren.svelte';
	import ViewDetailMushola from '$lib/components/views/ViewDetailMushola.svelte';
	import ViewDetailMT from '$lib/components/views/ViewDetailMT.svelte';

	let { data } = $props();

	// State management
	let situs = $state(null);
	let isLoading = $state(true);
	let error = $state(null);

	// Verification form state
	let verifyStatus = $state('');
	let verifySitusId = $state('');
	let isVerifying = $state(false);

	// Toast state
	let showToast = $state(false);
	let toastMessage = $state('');
	let toastType = $state('success');

	// Map state
	let map;
	let marker;

	// Fetch data on mount
	onMount(async () => {
		try {
			isLoading = true;
			error = null;
			const response = await apiService.get(`/situs/${$page.params.id}`);
			if (!response.ok) {
				throw new Error('Gagal memuat data situs');
			}
			situs = await response.json();

			// Pre-fill verification form
			verifyStatus = situs.status_verifikasi || '';
			verifySitusId = situs.situs_id || '';
		} catch (err) {
			console.error('Failed to fetch situs:', err);
			error = err.message;
		} finally {
			isLoading = false;
		}

		// Cleanup on unmount
		return () => {
			if (map) {
				map.remove();
				map = null;
				marker = null;
			}
		};
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
				scrollWheelZoom: false,
				zoomControl: true,
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

	// Format date to Indonesian locale
	function formatDate(dateString) {
		if (!dateString) return '-';
		const date = new Date(dateString);
		return date.toLocaleDateString('id-ID', {
			day: 'numeric',
			month: 'long',
			year: 'numeric'
		});
	}

	// Submit verification
	async function submitVerifikasi(e) {
		e.preventDefault();
		if (isVerifying) return; // Prevent double submit

		isVerifying = true;

		try {
			const res = await apiService.patch('/situs/' + $page.params.id + '/verify', {
				status_verifikasi: verifyStatus,
				situs_id: verifySitusId
			});

			if (res.ok) {
				toastMessage = 'Verifikasi berhasil disimpan';
				toastType = 'success';
				showToast = true;

				setTimeout(() => {
					location.reload();
				}, 2000);
			} else {
				const errText = await res.text();
				toastMessage = 'Gagal menyimpan verifikasi: ' + errText;
				toastType = 'error';
				showToast = true;
				isVerifying = false;

				setTimeout(() => {
					showToast = false;
				}, 4000);
			}
		} catch (err) {
			console.error('Verification error:', err);
			toastMessage = 'Terjadi kesalahan sistem saat menyimpan verifikasi';
			toastType = 'error';
			showToast = true;
			isVerifying = false;

			setTimeout(() => {
				showToast = false;
			}, 4000);
		}
	}

	function goBack() {
		goto('/dashboard/situs');
	}

	function getStatusBadgeClass(status) {
		if (!status) return 'badge-ghost';
		switch (status.toLowerCase()) {
			case 'terverifikasi':
				return 'badge-success';
			case 'ditolak':
				return 'badge-error';
			case 'menunggu':
				return 'badge-warning';
			default:
				return 'badge-ghost';
		}
	}
</script>

<svelte:head>
	<title>{situs?.nama || 'Detail Situs'} - Dashboard</title>
</svelte:head>

<div class="relative mx-auto max-w-7xl p-4">
	<Toast
		show={showToast}
		message={toastMessage}
		type={toastType}
		onclose={() => (showToast = false)}
	/>

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
						<div class="mt-1 flex flex-wrap items-center gap-2">
							<span class="badge badge-lg badge-primary">{situs.jenis_tipologi}</span>
							<span class="badge badge-outline">{situs.jenis_situs}</span>
							{#if situs.status_verifikasi}
								<span class="badge {getStatusBadgeClass(situs.status_verifikasi)}">
									{situs.status_verifikasi}
								</span>
							{/if}
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
						<div class="sm:col-span-2 lg:col-span-3">
							<p class="mb-1 text-xs font-medium text-base-content/70">Nomor Telepon Pengurus</p>
							{#if situs.nomor_telpon_pengurus && situs.nomor_telpon_pengurus.length > 0}
								<div class="flex flex-wrap gap-2">
									{#each situs.nomor_telpon_pengurus as phone}
										<span class="badge badge-outline badge-lg font-medium badge-primary"
											>{phone}</span
										>
									{/each}
								</div>
							{:else}
								<p class="font-medium">-</p>
							{/if}
						</div>

						<div>
							<p class="text-xs font-medium text-base-content/70">Nomor Telepon Instansi</p>
							<p class="font-medium">{situs.nomor_telepon || '-'}</p>
						</div>
						<div>
							<p class="text-xs font-medium text-base-content/70">Email</p>
							<p class="font-medium">{situs.email || '-'}</p>
						</div>
						<div>
							<p class="text-xs font-medium text-base-content/70">Website</p>
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
							<p class="text-xs font-medium text-base-content/70">Tahun Berdiri</p>
							<p class="font-medium">{situs.tahun_berdiri || '-'}</p>
						</div>
						<div>
							<p class="text-xs font-medium text-base-content/70">Nomor Badan Hukum</p>
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
							<p class="text-xs font-medium text-base-content/70">Alamat Lengkap</p>
							<p class="font-medium">{situs.alamat_lengkap || '-'}</p>
						</div>
						<div>
							<p class="text-xs font-medium text-base-content/70">Desa/Kelurahan</p>
							<p class="font-medium">{situs.desa || '-'}</p>
						</div>
						<div>
							<p class="text-xs font-medium text-base-content/70">Kecamatan</p>
							<p class="font-medium">{situs.kecamatan || '-'}</p>
						</div>
						<div>
							<p class="text-xs font-medium text-base-content/70">Kabupaten/Kota</p>
							<p class="font-medium">{situs.kabupaten_kota || '-'}</p>
						</div>
						<div>
							<p class="text-xs font-medium text-base-content/70">Provinsi</p>
							<p class="font-medium">{situs.provinsi || '-'}</p>
						</div>
						<div>
							<p class="text-xs font-medium text-base-content/70">Koordinat</p>
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
							<p class="text-xs font-medium text-base-content/70">Luas Tanah</p>
							<p class="font-medium">{situs.luas_tanah ? `${situs.luas_tanah} m²` : '-'}</p>
						</div>
						<div>
							<p class="text-xs font-medium text-base-content/70">Luas Bangunan</p>
							<p class="font-medium">{situs.luas_bangunan ? `${situs.luas_bangunan} m²` : '-'}</p>
						</div>
						<div>
							<p class="text-xs font-medium text-base-content/70">Status Tanah</p>
							<p class="font-medium">{situs.status_tanah || '-'}</p>
						</div>
						<div>
							<p class="text-xs font-medium text-base-content/70">Daya Tampung Max</p>
							<p class="font-medium">
								{situs.daya_tampung_max ? `${situs.daya_tampung_max} orang` : '-'}
							</p>
						</div>
						<div>
							<p class="text-xs font-medium text-base-content/70">Nomor AIW</p>
							<p class="font-medium">{situs.nomor_aiw || '-'}</p>
						</div>
						<div>
							<p class="text-xs font-medium text-base-content/70">No. Sertifikat Wakaf</p>
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

			{#if hasAnyPermission(data.user.permissions, 'situs:verify')}
				<div class="card border-2 border-primary bg-primary/5 shadow-lg">
					<div class="card-body p-4 sm:p-6">
						<h2 class="mb-2 card-title flex items-center gap-2 text-lg text-primary">
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
									d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
								/>
							</svg>
							Panel Verifikasi Situs
						</h2>
						<p class="mb-4 text-sm text-base-content/70">
							Status saat ini:
							{#if situs.status_verifikasi}
								<span class="badge {getStatusBadgeClass(situs.status_verifikasi)} ml-2">
									{situs.status_verifikasi}
								</span>
							{:else}
								<span class="ml-2 text-base-content/50">Belum ada status</span>
							{/if}
						</p>

						<form onsubmit={submitVerifikasi} class="space-y-4">
							<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
								<div class="form-control">
									<label class="label" for="verifyStatus">
										<span class="label-text font-medium">Status Verifikasi</span>
										<span class="label-text-alt text-error">*</span>
									</label>
									<select
										id="verifyStatus"
										class="select-bordered select w-full"
										required
										bind:value={verifyStatus}
									>
										<option value="" disabled>Pilih status...</option>
										<option value="terverifikasi">Terverifikasi</option>
										<option value="ditolak">Ditolak</option>
									</select>
								</div>

								<div class="form-control">
									<label class="label" for="verifySitusId">
										<span class="label-text font-medium">Nomor Statistik Kemenag</span>
										<span class="label-text-alt">Opsional</span>
									</label>
									<input
										id="verifySitusId"
										type="text"
										class="input-bordered input w-full"
										placeholder="Contoh: 51.2.12.345678"
										bind:value={verifySitusId}
										disabled={verifyStatus === 'ditolak'}
									/>
									{#if verifyStatus === 'ditolak'}
										<label for="verifySitusId" class="label">
											<span class="label-text-alt text-error"
												>Tidak dapat diisi saat status ditolak</span
											>
										</label>
									{/if}
								</div>
							</div>

							<div class="flex justify-end gap-2">
								<button type="submit" class="btn btn-primary" disabled={isVerifying}>
									{#if isVerifying}
										<span class="loading loading-spinner"></span> Menyimpan...
									{:else}
										<svg
											xmlns="http://www.w3.org/2000/svg"
											class="h-5 w-5"
											fill="none"
											viewBox="0 0 24 24"
											stroke="currentColor"
										>
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												stroke-width="2"
												d="M5 13l4 4L19 7"
											/>
										</svg>
										Simpan Verifikasi
									{/if}
								</button>
							</div>
						</form>
					</div>
				</div>
			{/if}
		</div>
	{/if}
</div>

<style>
	:global(.leaflet-container) {
		z-index: 0;
	}
</style>
