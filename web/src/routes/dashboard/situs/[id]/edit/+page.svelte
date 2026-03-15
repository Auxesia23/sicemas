<script>
	import { onMount, onDestroy } from 'svelte';
	import 'leaflet/dist/leaflet.css';
	import apiService from '$lib/api';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import DetailMasjidForm from '$lib/components/forms/DetailMasjidForm.svelte';
	import DetailPesantrenForm from '$lib/components/forms/DetailPesantrenForm.svelte';
	import DetailMusholaForm from '$lib/components/forms/DetailMusholaForm.svelte';
	import DetailMTForm from '$lib/components/forms/DetailMTForm.svelte';
	import Toast from '$lib/components/ui/Toast.svelte';
	import ConfirmModal from '$lib/components/ui/ConfirmModal.svelte';

	let { data } = $props();
	let user = $derived(data.user);

	// State for form data (no jenis_situs_id anymore)
	let formData = $state({
		nama: '',
		jenis_tipologi: '',
		nomor_telepon: '',
		email: '',
		website: '',
		nomor_badan_hukum: '',
		tahun_berdiri: '',
		provinsi: 'Jawa Barat',
		kabupaten_kota: 'Kab. Sukabumi',
		kecamatan: 'Ciemas',
		desa: '',
		alamat_lengkap: '',
		latitude: null,
		longitude: null,
		luas_tanah: null,
		luas_bangunan: null,
		status_tanah: '',
		nomor_aiw: '',
		nomor_sertifikat_wakaf: '',
		daya_tampung_max: null,
		detail: {}
	});

	// Simple state for jenis situs name (direct string from backend)
	let jenisSitusName = $state('');

	// State for images
	let existingPhotos = $state([]);
	let photosToDelete = $state([]);
	let newPhotos = $state([]);
	let isConverting = $state(false);
	let fileInputRef = $state();

	// State for map
	let map;
	let marker;

	// State for toast and modal
	let showToast = $state(false);
	let toastMessage = $state('');
	let toastType = $state('success');
	let showConfirmModal = $state(false);
	let isSubmitting = $state(false);
	let isLoading = $state(true);

	// Derived states based on direct string
	let isMasjid = $derived(jenisSitusName.toLowerCase().trim().includes('masjid'));
	let isPesantren = $derived(
		jenisSitusName.toLowerCase().trim().includes('ponpes') ||
			jenisSitusName.toLowerCase().trim().includes('pesantren')
	);
	let isMusholla = $derived(
		jenisSitusName.toLowerCase().trim().includes('mushola') ||
			jenisSitusName.toLowerCase().trim().includes('musholla')
	);
	let isMT = $derived(
		jenisSitusName.toLowerCase().trim() === 'mt' ||
			jenisSitusName.toLowerCase().trim().includes('majelis')
	);

	// Tipologi options
	const tipologiOptions = {
		Masjid: ['Masjid Besar', 'Masjid Jami', 'Masjid Bersejarah', 'Masjid Publik'],
		Mushola: ['Musholla'],
		Ponpes: ['Salafiyah', 'Kholafiyah', 'Salafiyah Wajar Dikdas', 'Muadalah'],
		MT: ['Majelis Taklim'],
		Madrasah: ['Madrasah Diniyah', 'Madrasah Ibtidaiyah', 'Madrasah Tsanawiyah', 'Madrasah Aliyah'],
		Lainnya: ['Umum', 'Khusus']
	};

	function getTipologiOptions() {
		if (!jenisSitusName) return [];
		const nameStr = jenisSitusName.toLowerCase().trim();

		if (nameStr === 'mt' || nameStr.includes('majelis')) return tipologiOptions['MT'];
		if (nameStr.includes('masjid')) return tipologiOptions['Masjid'];
		if (nameStr.includes('mushola') || nameStr.includes('musholla'))
			return tipologiOptions['Mushola'];
		if (nameStr.includes('ponpes') || nameStr.includes('pesantren'))
			return tipologiOptions['Ponpes'];
		if (nameStr.includes('madrasah')) return tipologiOptions['Madrasah'];

		return ['Umum', 'Khusus'];
	}

	// Map Action using Svelte action pattern
	function mapAction(node) {
		let L;

		// Gunakan setTimeout agar kontainer DOM selesai dirender sebelum peta diinisialisasi
		setTimeout(async () => {
			L = await import('leaflet');
			await import('leaflet/dist/leaflet.css');

			// Default coordinate
			const initialLat = formData.latitude || -6.2088;
			const initialLng = formData.longitude || 106.8456;

			map = L.map(node, {
				dragging: true,
				touchZoom: true,
				scrollWheelZoom: false,
				doubleClickZoom: false,
				boxZoom: false,
				keyboard: false,
				zoomControl: true,
				attributionControl: true
			}).setView([initialLat, initialLng], 15);

			L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
				attribution: '© OpenStreetMap contributors'
			}).addTo(map);

			// FIX BUGS VITE: Perbaiki link icon marker default Leaflet yang pecah
			delete L.Icon.Default.prototype._getIconUrl;
			L.Icon.Default.mergeOptions({
				iconRetinaUrl:
					'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon-2x.png',
				iconUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon.png',
				shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-shadow.png'
			});

			marker = L.marker([initialLat, initialLng], { draggable: true }).addTo(map);

			// FIX BUGS LEAFLET: Paksa peta menghitung ulang ukuran pembungkusnya
			setTimeout(() => {
				if (map) map.invalidateSize();
			}, 300);

			marker.on('dragend', (e) => {
				const latlng = e.target.getLatLng();
				formData.latitude = latlng.lat;
				formData.longitude = latlng.lng;
			});

			map.on('click', (e) => {
				formData.latitude = e.latlng.lat;
				formData.longitude = e.latlng.lng;
				marker.setLatLng(e.latlng);
			});
		}, 100);

		return {
			update(newCoords) {
				if (
					map &&
					marker &&
					newCoords &&
					newCoords.latitude !== null &&
					newCoords.longitude !== null
				) {
					const newLatLng = new L.LatLng(newCoords.latitude, newCoords.longitude);
					if (!marker.getLatLng().equals(newLatLng)) {
						marker.setLatLng(newLatLng);
						map.setView(newLatLng, 15);
					}
				}
			},
			destroy() {
				if (map) {
					map.remove();
					map = null;
					marker = null;
				}
			}
		};
	}

	onMount(async () => {
		isLoading = true;
		try {
			const response = await apiService.get(`/situs/${$page.params.id}`);
			if (!response.ok) {
				throw new Error('Gagal memuat data situs');
			}
			const siteData = await response.json();

			// Set string jenis situs
			jenisSitusName = siteData.jenis_situs || '';

			// Populate form data
			formData = {
				nama: siteData.nama || '',
				jenis_tipologi: siteData.jenis_tipologi || '',
				nomor_telepon: siteData.nomor_telepon || '',
				email: siteData.email || '',
				website: siteData.website || '',
				nomor_badan_hukum: siteData.nomor_badan_hukum || '',
				tahun_berdiri: siteData.tahun_berdiri ? Number(siteData.tahun_berdiri) : '',
				provinsi: siteData.provinsi || 'Jawa Barat',
				kabupaten_kota: siteData.kabupaten_kota || 'Kab. Sukabumi',
				kecamatan: siteData.kecamatan || 'Ciemas',
				desa: siteData.desa || '',
				alamat_lengkap: siteData.alamat_lengkap || '',
				latitude: siteData.latitude ?? null,
				longitude: siteData.longitude ?? null,
				luas_tanah: siteData.luas_tanah ?? null,
				luas_bangunan: siteData.luas_bangunan ?? null,
				status_tanah: siteData.status_tanah || '',
				nomor_aiw: siteData.nomor_aiw || '',
				nomor_sertifikat_wakaf: siteData.nomor_sertifikat_wakaf || '',
				daya_tampung_max: siteData.daya_tampung_max ?? null,
				detail: siteData.detail || {}
			};

			// Load foto existing
			existingPhotos = (siteData.foto || []).map((f) => ({
				id: f.id,
				url: f.image_url || f.url
			}));
		} catch (error) {
			console.error('Failed to fetch data:', error);
			toastMessage = 'Gagal memuat data situs';
			toastType = 'error';
			showToast = true;
		} finally {
			isLoading = false;
		}
	});

	onDestroy(() => {
		newPhotos.forEach((img) => {
			if (img?.preview) URL.revokeObjectURL(img.preview);
		});
		if (map) {
			map.remove();
			map = null;
		}
	});

	async function getCurrentLocation() {
		if (!navigator.geolocation) {
			toastMessage = 'Geolocation tidak didukung';
			toastType = 'error';
			showToast = true;
			return;
		}
		try {
			const position = await new Promise((resolve, reject) => {
				navigator.geolocation.getCurrentPosition(resolve, reject, {
					enableHighAccuracy: true,
					timeout: 15000
				});
			});
			formData.latitude = position.coords.latitude;
			formData.longitude = position.coords.longitude;
		} catch (error) {
			toastMessage = 'Gagal mendapatkan lokasi';
			toastType = 'error';
			showToast = true;
		}
	}

	async function convertToWebP(file) {
		return new Promise((resolve, reject) => {
			const reader = new FileReader();
			reader.onload = (e) => {
				const img = new Image();
				img.onload = () => {
					const canvas = document.createElement('canvas');
					canvas.width = img.width;
					canvas.height = img.height;
					const ctx = canvas.getContext('2d');
					ctx.drawImage(img, 0, 0);
					canvas.toBlob(
						(blob) => (blob ? resolve(blob) : reject(new Error('Failed'))),
						'image/webp',
						0.8
					);
				};
				img.onerror = () => reject(new Error('Failed to load'));
				img.src = e.target.result;
			};
			reader.onerror = () => reject(new Error('Failed to read'));
			reader.readAsDataURL(file);
		});
	}

	async function handlePhotoSelect(event) {
		const files = Array.from(event.target.files);
		if (files.length === 0) return;
		isConverting = true;
		try {
			for (const file of files) {
				if (!file.type.startsWith('image/')) continue;
				const preview = URL.createObjectURL(file);
				const webpBlob = await convertToWebP(file);
				newPhotos.push({ file, preview, webpBlob });
			}
		} catch (error) {
			toastMessage = 'Gagal mengkonversi gambar';
			toastType = 'error';
			showToast = true;
		} finally {
			isConverting = false;
			if (fileInputRef) fileInputRef.value = '';
		}
	}

	function removeNewPhoto(index) {
		if (newPhotos[index]?.preview) URL.revokeObjectURL(newPhotos[index].preview);
		newPhotos.splice(index, 1);
	}

	function markPhotoForDeletion(photoId) {
		photosToDelete.push(photoId);
		existingPhotos = existingPhotos.filter((p) => p.id !== photoId);
	}

	async function uploadNewImages(situsId) {
		if (newPhotos.length === 0) return true;
		const formDataImages = new FormData();
		for (const img of newPhotos) {
			formDataImages.append('images', img.webpBlob, `image-${Date.now()}.webp`);
		}
		const response = await apiService.post(`/situs/${situsId}/foto`, formDataImages);
		if (!response.ok) throw new Error('Gagal mengunggah foto baru');
		return true;
	}

	async function deleteMarkedImages(situsId) {
		if (photosToDelete.length === 0) return true;
		const response = await apiService.delete(`/situs/${situsId}/foto`, { ids: photosToDelete });
		if (!response.ok) throw new Error('Gagal menghapus foto lama');
		return true;
	}

	function confirmUpdate() {
		showConfirmModal = true;
	}

	async function handleUpdate() {
		isSubmitting = true;
		showConfirmModal = false;

		try {
			// Parsing untuk Zod di Backend Go
			formData.tahun_berdiri = parseInt(formData.tahun_berdiri) || 0;
			formData.daya_tampung_max = parseInt(formData.daya_tampung_max) || 0;
			formData.luas_tanah = parseFloat(formData.luas_tanah) || 0;
			formData.luas_bangunan = parseFloat(formData.luas_bangunan) || 0;

			// STEP A: PUT data text (must complete first)
			const updateResponse = await apiService.put(`/situs/${$page.params.id}`, formData);
			if (!updateResponse.ok) {
				const errorText = await updateResponse.text();
				throw new Error(errorText || 'Gagal memperbarui data situs');
			}

			const situsId = $page.params.id;

			// STEP B & C: Parallel photo processing
			const uploadPromise = newPhotos.length > 0 ? uploadNewImages(situsId) : Promise.resolve(true);
			const deletePromise =
				photosToDelete.length > 0 ? deleteMarkedImages(situsId) : Promise.resolve(true);

			await Promise.all([uploadPromise, deletePromise]);

			toastMessage = 'Data situs berhasil diperbarui';
			toastType = 'success';
			showToast = true;

			newPhotos.forEach((img) => {
				if (img?.preview) URL.revokeObjectURL(img.preview);
			});

			setTimeout(() => {
				showToast = false;
				goto('/dashboard/situs');
			}, 2000);
		} catch (error) {
			console.error('Update error:', error);
			toastMessage = error.message || 'Terjadi kesalahan saat memperbarui data';
			toastType = 'error';
			showToast = true;
		} finally {
			isSubmitting = false;
		}
	}

	function goBack() {
		goto('/dashboard/situs');
	}
</script>

<svelte:head>
	<title>Edit {formData.nama || 'Situs'} - Dashboard</title>
</svelte:head>

<div class="relative min-h-screen bg-base-200 p-3 sm:p-4">
	<Toast
		show={showToast}
		message={toastMessage}
		type={toastType}
		onclose={() => (showToast = false)}
	/>

	<ConfirmModal
		show={showConfirmModal}
		title="Konfirmasi Pembaruan"
		message="Apakah Anda yakin ingin memperbarui data situs ini?"
		confirmText="Simpan"
		cancelText="Batal"
		isProcessing={isSubmitting}
		onConfirm={handleUpdate}
		onCancel={() => (showConfirmModal = false)}
	/>

	{#if isLoading}
		<div class="flex h-[80vh] items-center justify-center">
			<span class="loading loading-lg loading-spinner text-primary"></span>
		</div>
	{:else}
		<div class="mx-auto max-w-3xl">
			<div class="mb-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
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
						<h1 class="text-2xl font-bold text-primary sm:text-3xl">Edit Situs Keagamaan</h1>
						<p class="text-sm text-base-content/70">Perbarui informasi situs keagamaan</p>
					</div>
				</div>
			</div>

			<form
				class="space-y-6"
				onsubmit={(e) => {
					e.preventDefault();
					confirmUpdate();
				}}
			>
				<div class="card border border-base-200 bg-base-100 shadow-lg">
					<div class="card-body p-4 sm:p-6">
						<h2 class="mb-4 card-title text-lg">Jenis Situs</h2>
						<div class="form-control w-full">
							<label class="label" for="jenis_situs_display">
								<span class="label-text font-medium">Jenis Situs (Tidak dapat diubah)</span>
							</label>
							<input
								id="jenis_situs_display"
								type="text"
								class="input-bordered input min-h-11 w-full bg-base-200 font-semibold text-primary"
								value={jenisSitusName}
								readonly
							/>
						</div>
					</div>
				</div>

				<div class="card border border-base-200 bg-base-100 shadow-lg">
					<div class="card-body p-4 sm:p-6">
						<h2 class="mb-4 card-title text-lg">Informasi Umum</h2>
						<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
							<div class="form-control">
								<label class="label" for="nama">
									<span class="label-text font-medium"
										>Nama Situs <span class="text-error">*</span></span
									>
								</label>
								<input
									required
									id="nama"
									type="text"
									class="input-bordered input min-h-11 w-full"
									bind:value={formData.nama}
								/>
							</div>
							<div class="form-control">
								<label class="label" for="jenis_tipologi">
									<span class="label-text font-medium"
										>Jenis Tipologi <span class="text-error">*</span></span
									>
								</label>
								{#if getTipologiOptions().length <= 1}
									<input
										id="jenis_tipologi"
										type="text"
										class="input-bordered input min-h-11 w-full bg-base-200"
										bind:value={formData.jenis_tipologi}
										readonly
									/>
								{:else}
									<select
										id="jenis_tipologi"
										class="select-bordered select min-h-11 w-full"
										required
										bind:value={formData.jenis_tipologi}
									>
										<option value="" disabled>Pilih jenis tipologi...</option>
										{#each getTipologiOptions() as option}
											<option value={option}>{option}</option>
										{/each}
									</select>
								{/if}
							</div>
							<div class="form-control">
								<label class="label" for="nomor_telepon"
									><span class="label-text font-medium">Nomor Telepon</span></label
								>
								<input
									id="nomor_telepon"
									type="tel"
									class="input-bordered input min-h-11 w-full"
									bind:value={formData.nomor_telepon}
								/>
							</div>
							<div class="form-control">
								<label class="label" for="email"
									><span class="label-text font-medium">Email</span></label
								>
								<input
									id="email"
									type="email"
									class="input-bordered input min-h-11 w-full"
									bind:value={formData.email}
								/>
							</div>
							<div class="form-control">
								<label class="label" for="website"
									><span class="label-text font-medium">Website</span></label
								>
								<input
									id="website"
									type="url"
									class="input-bordered input min-h-11 w-full"
									bind:value={formData.website}
								/>
							</div>
							<div class="form-control">
								<label class="label" for="nomor_badan_hukum"
									><span class="label-text font-medium">Nomor Badan Hukum</span></label
								>
								<input
									id="nomor_badan_hukum"
									type="text"
									class="input-bordered input min-h-11 w-full"
									bind:value={formData.nomor_badan_hukum}
								/>
							</div>
							<div class="form-control">
								<label class="label" for="tahun_berdiri"
									><span class="label-text font-medium"
										>Tahun Berdiri <span class="text-error">*</span></span
									></label
								>
								<input
									required
									id="tahun_berdiri"
									type="number"
									min="0"
									class="input-bordered input min-h-11 w-full"
									bind:value={formData.tahun_berdiri}
								/>
							</div>
						</div>
					</div>
				</div>

				<div class="card border border-base-200 bg-base-100 shadow-lg">
					<div class="card-body p-4 sm:p-6">
						<h2 class="mb-4 card-title text-lg">Lokasi</h2>
						<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
							<div class="form-control">
								<label class="label" for="provinsi"
									><span class="label-text font-medium">Provinsi</span></label
								>
								<input
									id="provinsi"
									type="text"
									class="input-bordered input min-h-11 w-full bg-base-200"
									value={formData.provinsi}
									readonly
								/>
							</div>
							<div class="form-control">
								<label class="label" for="kabupaten_kota"
									><span class="label-text font-medium">Kabupaten/Kota</span></label
								>
								<input
									id="kabupaten_kota"
									type="text"
									class="input-bordered input min-h-11 w-full bg-base-200"
									value={formData.kabupaten_kota}
									readonly
								/>
							</div>
							<div class="form-control">
								<label class="label" for="kecamatan"
									><span class="label-text font-medium">Kecamatan</span></label
								>
								<input
									id="kecamatan"
									type="text"
									class="input-bordered input min-h-11 w-full bg-base-200"
									value={formData.kecamatan}
									readonly
								/>
							</div>
							<div class="form-control">
								<label class="label" for="desa"
									><span class="label-text font-medium"
										>Desa/Kelurahan <span class="text-error">*</span></span
									></label
								>
								<input
									required
									id="desa"
									type="text"
									class="input-bordered input min-h-11 w-full"
									bind:value={formData.desa}
								/>
							</div>
							<div class="form-control sm:col-span-2">
								<label class="label" for="alamat_lengkap"
									><span class="label-text font-medium"
										>Alamat Lengkap <span class="text-error">*</span></span
									></label
								>
								<textarea
									required
									id="alamat_lengkap"
									class="textarea-bordered textarea min-h-20 w-full"
									bind:value={formData.alamat_lengkap}
									rows="2"
								></textarea>
							</div>
						</div>

						<div class="mt-4 sm:mt-6">
							<button
								type="button"
								class="btn mb-3 btn-outline sm:mb-4"
								onclick={getCurrentLocation}
							>
								Gunakan Lokasi Saat Ini
							</button>
							<div class="mb-3 grid grid-cols-2 gap-2 sm:mb-4 sm:gap-4">
								<div class="form-control">
									<label class="label" for="latitude"
										><span class="label-text font-medium"
											>Latitude <span class="text-error">*</span></span
										></label
									>
									<input
										required
										id="latitude"
										type="number"
										step="any"
										class="input-bordered input min-h-11 w-full"
										bind:value={formData.latitude}
									/>
								</div>
								<div class="form-control">
									<label class="label" for="longitude"
										><span class="label-text font-medium"
											>Longitude <span class="text-error">*</span></span
										></label
									>
									<input
										required
										id="longitude"
										type="number"
										step="any"
										class="input-bordered input min-h-11 w-full"
										bind:value={formData.longitude}
									/>
								</div>
							</div>
							<div
								use:mapAction={formData}
								class="relative z-0 h-64 w-full overflow-hidden rounded-lg border border-base-300 sm:h-80"
							></div>
							<p class="mt-2 text-xs text-base-content/70">
								Klik pada peta atau geser marker untuk mengubah lokasi
							</p>
						</div>
					</div>
				</div>

				<div class="card border border-base-200 bg-base-100 shadow-lg">
					<div class="card-body p-4 sm:p-6">
						<h2 class="mb-4 card-title text-lg">Kapasitas & Legalitas</h2>
						<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
							<div class="form-control">
								<label class="label" for="luas_tanah"
									><span class="label-text font-medium"
										>Luas Tanah (m²) <span class="text-error">*</span></span
									></label
								>
								<input
									required
									id="luas_tanah"
									type="number"
									min="0"
									step="any"
									class="input-bordered input min-h-11 w-full"
									bind:value={formData.luas_tanah}
								/>
							</div>
							<div class="form-control">
								<label class="label" for="luas_bangunan"
									><span class="label-text font-medium"
										>Luas Bangunan (m²) <span class="text-error">*</span></span
									></label
								>
								<input
									required
									id="luas_bangunan"
									type="number"
									min="0"
									step="any"
									class="input-bordered input min-h-11 w-full"
									bind:value={formData.luas_bangunan}
								/>
							</div>
							<div class="form-control">
								<label class="label" for="status_tanah"
									><span class="label-text font-medium"
										>Status Tanah <span class="text-error">*</span></span
									></label
								>
								<select
									required
									id="status_tanah"
									class="select-bordered select min-h-11 w-full"
									bind:value={formData.status_tanah}
								>
									<option value="" disabled>Pilih status...</option>
									<option value="Wakaf">Wakaf</option>
									<option value="Milik Sendiri">Milik Sendiri</option>
									<option value="Sewa">Sewa</option>
									<option value="Hibah">Hibah</option>
									<option value="Lainnya">Lainnya</option>
								</select>
							</div>
							<div class="form-control">
								<label class="label" for="daya_tampung_max"
									><span class="label-text font-medium"
										>Daya Tampung Max <span class="text-error">*</span></span
									></label
								>
								<input
									required
									id="daya_tampung_max"
									type="number"
									min="0"
									class="input-bordered input min-h-11 w-full"
									bind:value={formData.daya_tampung_max}
								/>
							</div>
							<div class="form-control">
								<label class="label" for="nomor_aiw"
									><span class="label-text font-medium">Nomor AIW</span></label
								>
								<input
									id="nomor_aiw"
									type="text"
									class="input-bordered input min-h-11 w-full"
									bind:value={formData.nomor_aiw}
								/>
							</div>
							<div class="form-control">
								<label class="label" for="nomor_sertifikat_wakaf"
									><span class="label-text font-medium">No. Sertifikat Wakaf</span></label
								>
								<input
									id="nomor_sertifikat_wakaf"
									type="text"
									class="input-bordered input min-h-11 w-full"
									bind:value={formData.nomor_sertifikat_wakaf}
								/>
							</div>
						</div>
					</div>
				</div>

				{#if isMasjid}
					<DetailMasjidForm bind:detail={formData.detail} />
				{/if}
				{#if isPesantren}
					<DetailPesantrenForm bind:detail={formData.detail} />
				{/if}
				{#if isMusholla}
					<DetailMusholaForm bind:detail={formData.detail} />
				{/if}
				{#if isMT}
					<DetailMTForm bind:detail={formData.detail} />
				{/if}

				<div class="card border border-base-200 bg-base-100 shadow-lg">
					<div class="card-body p-4 sm:p-6">
						<h2 class="mb-4 card-title text-lg">Foto Situs</h2>

						{#if existingPhotos.length > 0}
							<div class="mb-4">
								<label for="existing-photos" class="label"
									><span class="label-text font-medium">Foto Tersimpan</span></label
								>
								<div class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4">
									{#each existingPhotos as photo (photo.id)}
										<div class="relative overflow-hidden rounded-lg border border-base-300">
											<img src={photo.url} alt="Foto situs" class="h-32 w-full object-cover" />
											<button
												type="button"
												class="btn absolute top-1 right-1 btn-circle btn-xs btn-error"
												onclick={() => markPhotoForDeletion(photo.id)}
												title="Hapus foto ini">✕</button
											>
										</div>
									{/each}
								</div>
							</div>
						{/if}

						{#if newPhotos.length > 0}
							<div class="mb-4">
								<label for="new-photos" class="label"
									><span class="label-text font-medium">Foto Baru (akan diunggah)</span></label
								>
								<div class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4">
									{#each newPhotos as photo, index (photo.preview)}
										<div class="relative overflow-hidden rounded-lg border border-base-300">
											<img src={photo.preview} alt="Preview" class="h-32 w-full object-cover" />
											<button
												type="button"
												class="btn absolute top-1 right-1 btn-circle btn-xs btn-error"
												onclick={() => removeNewPhoto(index)}
												title="Batalkan unggah">✕</button
											>
											<div
												class="absolute bottom-1 left-1 rounded bg-black/70 px-1 text-xs text-white"
											>
												WebP
											</div>
										</div>
									{/each}
								</div>
							</div>
						{/if}

						<div class="form-control">
							<label class="label" for="images"
								><span class="label-text font-medium">Unggah Foto Tambahan</span></label
							>
							<input
								id="images"
								type="file"
								accept="image/*"
								multiple
								class="file-input-bordered file-input min-h-11 w-full"
								onchange={handlePhotoSelect}
								disabled={isConverting}
								bind:this={fileInputRef}
							/>
							{#if isConverting}
								<div class="mt-2 flex items-center gap-2 text-sm text-base-content/70">
									<span class="loading loading-sm loading-spinner"></span>
									Mengkonversi gambar ke WebP...
								</div>
							{/if}
						</div>
					</div>
				</div>

				<div class="flex flex-col justify-end gap-3 pb-4 sm:flex-row sm:pb-8">
					<button
						type="button"
						class="btn min-h-11 w-full btn-ghost sm:w-auto"
						onclick={goBack}
						disabled={isSubmitting}>Batal</button
					>
					<button
						type="submit"
						class="btn min-h-11 w-full btn-primary sm:w-auto"
						disabled={isSubmitting}
					>
						{#if isSubmitting}
							<span class="loading loading-spinner"></span> Menyimpan...
						{:else}
							Simpan Perubahan
						{/if}
					</button>
				</div>
			</form>
		</div>
	{/if}
</div>

<style>
	/* Prevent Leaflet map from overlapping dropdowns/modals */
	:global(.leaflet-container) {
		z-index: 10;
	}
	:global(.leaflet-control-container) {
		z-index: 10;
	}
</style>
