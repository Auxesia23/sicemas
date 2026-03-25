<script>
	import { onMount, onDestroy } from 'svelte';
	import 'leaflet/dist/leaflet.css';
	import apiService from '$lib/api';
	import { goto } from '$app/navigation';
	import DetailMasjidForm from '$lib/components/forms/DetailMasjidForm.svelte';
	import DetailPesantrenForm from '$lib/components/forms/DetailPesantrenForm.svelte';
	import DetailMusholaForm from '$lib/components/forms/DetailMusholaForm.svelte';
	import DetailMTForm from '$lib/components/forms/DetailMTForm.svelte';
	import Toast from '$lib/components/ui/Toast.svelte';

	import { z } from 'zod';

	let formData = $state({
		jenis_situs_id: '',
		nama: '',
		jenis_tipologi: '',
		nomor_telepon: '',
		nomor_telpon_pengurus: [],
		email: '',
		website: '',
		nomor_badan_hukum: '',
		tahun_berdiri: 2024, // Default tahun
		provinsi: 'Jawa Barat',
		kabupaten_kota: 'Kab. Sukabumi',
		kecamatan: 'Ciemas',
		desa: '',
		alamat_lengkap: '',
		latitude: null, // Khusus koordinat biarkan null dulu biar divalidasi peta
		longitude: null,
		luas_tanah: 0, // Default 0
		luas_bangunan: 0, // Default 0
		status_tanah: '',
		nomor_aiw: '',
		nomor_sertifikat_wakaf: '',
		daya_tampung_max: 0, // Default 0
		detail: {}
	});
	let validationErrors = $state({});

	const detailPesantrenSchema = z.object({
		nama_yayasan: z.string().optional(),
		pimpinan_pesantren: z.string().min(1, 'Nama Pimpinan wajib diisi'),
		kepengurusan: z.object({
			ketua: z.string().min(1, 'Ketua wajib diisi'),
			sekretaris: z.array(z.string()).min(1, 'Minimal 1 sekretaris wajib'),
			bendahara: z.array(z.string()).min(1, 'Minimal 1 bendahara wajib')
		}),
		santri: z.object({
			mondok: z.object({ pria: z.number().min(0), wanita: z.number().min(0) }),
			tidak_mondok: z.object({ pria: z.number().min(0), wanita: z.number().min(0) }),
			total: z.object({ pria: z.number().min(0), wanita: z.number().min(0) })
		})
	});

	const detailMTSchema = z.object({
		pengurus: z.object({
			ketua: z.string().min(1, 'Ketua wajib diisi'),
			sekretaris: z.array(z.string()).min(1, 'Minimal 1 sekretaris'),
			bendahara: z.array(z.string()).min(1, 'Minimal 1 bendahara')
		}),
		penceramah: z
			.array(
				z.object({
					nama: z.string().min(1),
					materi: z.string().min(1)
				})
			)
			.min(1, 'Minimal 1 penceramah wajib'),
		kehadiran_pria: z
			.object({
				jumlah_dewasa: z.number().min(0),
				jumlah_remaja: z.number().min(0),
				waktu_pengajian: z.string()
			})
			.optional(),
		kehadiran_wanita: z
			.object({
				jumlah_dewasa: z.number().min(0),
				jumlah_remaja: z.number().min(0),
				waktu_pengajian: z.string()
			})
			.optional()
	});

	const detailMusholaSchema = z.object({
		sdm_masjid: z.object({
			jumlah_imam: z.number().min(1, 'Jumlah imam wajib diisi'),
			jumlah_muadzin: z.number().min(1, 'Jumlah muadzin wajib diisi'),
			jumlah_jemaah: z.number().min(1, 'Jumlah jemaah wajib diisi')
		}),
		nama_imam: z.array(z.string()).min(1, 'Minimal 1 nama imam wajib diisi'),
		nama_muadzin: z.array(z.string()).min(1, 'Minimal 1 nama muadzin wajib diisi')
	});

	const detailMasjidSchema = z.object({
		sdm_masjid: z.object({
			jumlah_pengurus: z.number().min(1, 'Jumlah pengurus wajib diisi'),
			jumlah_imam: z.number().min(1, 'Jumlah imam wajib diisi'),
			jumlah_khotib: z.number().min(1, 'Jumlah khotib wajib diisi'),
			jumlah_muadzin: z.number().min(1, 'Jumlah muadzin wajib diisi'),
			jumlah_remaja: z.number().min(1, 'Jumlah remaja wajib diisi'),
			jumlah_jemaah: z.number().min(1, 'Jumlah jemaah wajib diisi')
		}),
		pengurus_dkm: z.object({
			ketua: z.string().min(1, 'Ketua DKM wajib diisi'),
			sekretaris: z.array(z.string()).min(1, 'Minimal 1 sekretaris wajib diisi'),
			bendahara: z.array(z.string()).min(1, 'Minimal 1 bendahara wajib diisi'),
			nama_imam: z.array(z.string()).min(1, 'Minimal 1 nama imam wajib diisi'),
			nama_muazdin: z.array(z.string()).min(1, 'Minimal 1 nama muadzin wajib diisi')
		})
	});

	const situsSchema = z.object({
		jenis_situs_id: z.string().min(1, 'Jenis situs wajib dipilih'),
		nama: z.string().min(1, 'Nama situs wajib diisi'),
		jenis_tipologi: z.string().min(1, 'Tipologi wajib dipilih'),
		tahun_berdiri: z.number().min(1000, 'Tahun tidak valid'),
		provinsi: z.string().min(1),
		kabupaten_kota: z.string().min(1),
		kecamatan: z.string().min(1),
		desa: z.string().min(1, 'Desa wajib diisi'),
		alamat_lengkap: z.string().min(5, 'Alamat kurang lengkap'),
		latitude: z.number({ invalid_type_error: 'Tentukan lokasi di peta' }),
		longitude: z.number({ invalid_type_error: 'Tentukan lokasi di peta' }),
		luas_tanah: z.number().min(1, 'Luas tanah wajib diisi'),
		luas_bangunan: z.number().min(1, 'Luas bangunan wajib diisi'),
		nomor_telpon_pengurus: z.array(z.string()).min(1, 'Minimal 1 nomor pengurus'),
		status_tanah: z.string().min(1, 'Status tanah wajib dipilih'),
		daya_tampung_max: z.number().min(1, 'Daya tampung wajib diisi')
	});

	let currentPhoneInput = $state('');

	function addPhone() {
		const phone = currentPhoneInput.trim();
		if (phone && !formData.nomor_telpon_pengurus.includes(phone)) {
			formData.nomor_telpon_pengurus = [...formData.nomor_telpon_pengurus, phone];
			currentPhoneInput = '';
			if (validationErrors.nomor_telpon_pengurus) validationErrors.nomor_telpon_pengurus = null;
		}
	}

	function removePhone(index) {
		formData.nomor_telpon_pengurus = formData.nomor_telpon_pengurus.filter((_, i) => i !== index);
	}

	const tipologiOptions = {
		Masjid: ['Masjid Besar', 'Masjid Jami', 'Masjid Bersejarah', 'Masjid Publik'],
		Mushola: ['Mushola'],
		Ponpes: ['Salafiyah', 'Kholafiyah', 'Salafiyah Wajar Dikdas', 'Muadalah'],
		MT: ['Majelis Taklim'],
		Madrasah: ['Madrasah Diniyah', 'Madrasah Ibtidaiyah', 'Madrasah Tsanawiyah', 'Madrasah Aliyah'],
		Lainnya: ['Umum', 'Khusus']
	};

	function getTipologiOptions() {
		if (!selectedJenisSitusName) return [];
		for (const [key, options] of Object.entries(tipologiOptions)) {
			if (selectedJenisSitusName.toLowerCase().includes(key.toLowerCase())) return options;
		}
		return ['Umum', 'Khusus'];
	}

	let jenisSitusList = $state([]);
	let selectedJenisSitusName = $state('');
	let mapContainer;
	let map;
	let marker;
	let showToast = $state(false);
	let toastMessage = $state('');
	let toastType = $state('success');
	let selectedImages = $state([]);
	let isConverting = $state(false);
	let fileInputRef;
	let isSubmitting = $state(false);

	// --- FIX BUGS MATCHING STRING ---
	let isMasjid = $derived(selectedJenisSitusName.toLowerCase().includes('masjid'));
	let isPesantren = $derived(
		selectedJenisSitusName.toLowerCase().includes('ponpes') ||
			selectedJenisSitusName.toLowerCase().includes('pesantren')
	);
	let isMusholla = $derived(
		selectedJenisSitusName.toLowerCase().includes('mushola') ||
			selectedJenisSitusName.toLowerCase().includes('musholla')
	);
	let isMT = $derived(
		selectedJenisSitusName.toLowerCase() === 'mt' ||
			selectedJenisSitusName.toLowerCase().includes('majelis')
	);

	onMount(async () => {
		try {
			const response = await apiService.get('/jenis-situs');
			const data = await response.json();
			jenisSitusList = data.data || data;
		} catch (error) {
			console.error('Failed to fetch jenis situs:', error);
		}
		if (mapContainer && !map) {
			const L = await import('leaflet');
			await import('leaflet/dist/leaflet.css');

			delete L.Icon.Default.prototype._getIconUrl;
			L.Icon.Default.mergeOptions({
				iconRetinaUrl:
					'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-icon-2x.png',
				iconUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-icon.png',
				shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-shadow.png'
			});

			map = L.map(mapContainer).setView([-6.2088, 106.8456], 13);
			L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
				attribution: '© OpenStreetMap contributors'
			}).addTo(map);
			marker = L.marker([-6.2088, 106.8456], { draggable: true }).addTo(map);

			marker.on('dragend', (e) => {
				const latlng = e.target.getLatLng();
				formData.latitude = latlng.lat;
				formData.longitude = latlng.lng;
				validationErrors.latitude = null;
			});

			map.on('click', (e) => {
				formData.latitude = e.latlng.lat;
				formData.longitude = e.latlng.lng;
				marker.setLatLng(e.latlng);
				validationErrors.latitude = null;
			});
		}
	});

	onDestroy(() => {
		selectedImages.forEach((img) => {
			if (img?.preview) URL.revokeObjectURL(img.preview);
		});
	});

	function handleJenisSitusChange(event) {
		const selectedId = event.target.value;
		formData.jenis_situs_id = selectedId;
		validationErrors.jenis_situs_id = null;

		const selected = jenisSitusList.find((item) => item.id === selectedId);
		if (selected) {
			selectedJenisSitusName = selected.nama || selected.name || '';
			formData.jenis_tipologi = '';

			const nameLower = selectedJenisSitusName.toLowerCase();

			// --- FIX BUGS MATCHING STRING ---
			if (nameLower.includes('masjid')) {
				formData.detail = {};
				formData.jenis_tipologi = tipologiOptions.Masjid[0];
			} else if (nameLower.includes('ponpes') || nameLower.includes('pesantren')) {
				formData.detail = {};
				formData.jenis_tipologi = tipologiOptions.Ponpes[0];
			} else if (nameLower.includes('mushola') || nameLower.includes('musholla')) {
				formData.detail = {};
				formData.jenis_tipologi = tipologiOptions.Mushola[0];
			} else if (nameLower === 'mt' || nameLower.includes('majelis')) {
				formData.detail = {};
				formData.jenis_tipologi = tipologiOptions.MT[0];
			} else {
				formData.detail = {};
			}
		}
	}

	async function getCurrentLocation() {
		if (!navigator.geolocation) {
			alert('Geolocation is not supported by your browser');
			return;
		}
		const btn = event.currentTarget;
		const originalText = btn.innerHTML;
		btn.innerHTML = '<span class="loading loading-spinner loading-sm"></span> Mencari...';
		btn.disabled = true;

		try {
			const position = await new Promise((resolve, reject) => {
				navigator.geolocation.getCurrentPosition(resolve, reject, {
					enableHighAccuracy: true,
					timeout: 15000,
					maximumAge: 0
				});
			});
			formData.latitude = position.coords.latitude;
			formData.longitude = position.coords.longitude;
			validationErrors.latitude = null;

			if (map && marker) {
				marker.setLatLng([formData.latitude, formData.longitude]);
				map.flyTo([formData.latitude, formData.longitude], 16, { animate: true, duration: 1.5 });
			}
		} catch (error) {
			alert('Gagal mendapatkan lokasi. Silakan pilih manual di peta.');
		} finally {
			btn.innerHTML = originalText;
			btn.disabled = false;
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
				img.onerror = () => reject(new Error('Failed'));
				img.src = e.target.result;
			};
			reader.readAsDataURL(file);
		});
	}

	async function handleImageSelect(event) {
		const files = Array.from(event.target.files);
		if (files.length === 0) return;
		isConverting = true;
		try {
			for (const file of files) {
				if (!file.type.startsWith('image/')) continue;
				const preview = URL.createObjectURL(file);
				const webpBlob = await convertToWebP(file);
				selectedImages.push({ file, preview, webpBlob });
			}
		} catch (error) {
			console.error(error);
		} finally {
			isConverting = false;
		}
	}

	function removeImage(index) {
		if (selectedImages[index]?.preview) URL.revokeObjectURL(selectedImages[index].preview);
		selectedImages.splice(index, 1);
		if (selectedImages.length === 0 && fileInputRef) fileInputRef.value = '';
	}

	async function uploadImages(situsId) {
		if (selectedImages.length === 0) return;
		const formDataImages = new FormData();
		for (const img of selectedImages)
			formDataImages.append('images', img.webpBlob, `image-${Date.now()}.webp`);
		await apiService.post(`/situs/${situsId}/foto`, formDataImages);
	}

	async function handleSubmit() {
		let errors = {};
		let hasError = false;

		const mapZodErrors = (issues) => {
			hasError = true;
			issues.forEach((issue) => {
				const path = issue.path;

				if (path.length === 1) {
					errors[path[0]] = issue.message;
				} else if (path.length === 2) {
					if (!errors[path[0]]) errors[path[0]] = {};
					errors[path[0]][path[1]] = issue.message;
				} else if (path.length === 3) {
					if (!errors[path[0]]) errors[path[0]] = {};
					if (!errors[path[0]][path[1]]) errors[path[0]][path[1]] = {};
					errors[path[0]][path[1]][path[2]] = issue.message;
				}
			});
		};

		const result = situsSchema.safeParse(formData);
		if (!result.success) mapZodErrors(result.error.issues);

		if (isMasjid) {
			const detailResult = detailMasjidSchema.safeParse(formData.detail);
			if (!detailResult.success) mapZodErrors(detailResult.error.issues);
		} else if (isMusholla) {
			const detailResult = detailMusholaSchema.safeParse(formData.detail);
			if (!detailResult.success) mapZodErrors(detailResult.error.issues);
		} else if (isMT) {
			const detailResult = detailMTSchema.safeParse(formData.detail);
			if (!detailResult.success) mapZodErrors(detailResult.error.issues);
		} else if (isPesantren) {
			const detailResult = detailPesantrenSchema.safeParse(formData.detail);
			if (!detailResult.success) mapZodErrors(detailResult.error.issues);
		}

		if (hasError) {
			validationErrors = errors;

			toastMessage = 'Tolong periksa kembali form isian Anda';
			toastType = 'error';
			showToast = true;
			setTimeout(() => (showToast = false), 3000);
			window.scrollTo({ top: 0, behavior: 'smooth' });
			return;
		}

		validationErrors = {};
		isSubmitting = true;

		try {
			const response = await apiService.post('/situs', formData);
			const data = await response.json();

			if (response.status === 201) {
				const situsId = data.id || data.data?.id;
				if (situsId && selectedImages.length > 0) {
					await uploadImages(situsId);
				}

				toastMessage = 'Data situs berhasil disimpan';
				toastType = 'success';
				showToast = true;

				setTimeout(() => {
					goto('/dashboard/situs');
				}, 2000);
			} else {
				throw new Error('Gagal simpan');
			}
		} catch (error) {
			console.error('Failed to submit form:', error);
			toastMessage = 'Terjadi kesalahan. Silakan coba lagi.';
			toastType = 'error';
			showToast = true;
			setTimeout(() => (showToast = false), 5000);
			isSubmitting = false;
		}
	}
</script>

<svelte:head>
	<title>Tambah Situs Keagamaan</title>
</svelte:head>

<div class="min-h-screen bg-base-200 p-3 sm:p-4">
	<div class="mx-auto max-w-3xl">
		<div class="mb-4 sm:mb-6">
			<h1 class="text-2xl font-bold text-primary sm:text-3xl">Pendataan Situs Keagamaan</h1>
			<p class="mt-1 text-sm text-base-content/70 sm:mt-2 sm:text-base">
				Formulir pendaftaran situs keagamaan baru
			</p>
		</div>

		<form
			onsubmit={(e) => {
				e.preventDefault();
				handleSubmit();
			}}
			class="space-y-3 sm:space-y-4"
		>
			<div class="card bg-base-100 shadow-md sm:shadow-xl">
				<div class="card-body p-4 sm:p-5">
					<h2 class="mb-3 card-title text-base sm:text-lg">Jenis Situs</h2>
					<div class="form-control w-full">
						<label class="label" for="jenis_situs_id">
							<span class="label-text font-medium"
								>Pilih Jenis Situs <span class="text-error">*</span></span
							>
						</label>
						<select
							id="jenis_situs_id"
							class="select-bordered select min-h-11 w-full {validationErrors.jenis_situs_id
								? 'select-error'
								: ''}"
							bind:value={formData.jenis_situs_id}
							onchange={handleJenisSitusChange}
						>
							<option value="" disabled>Pilih jenis situs...</option>
							{#each jenisSitusList as jenis (jenis.id)}
								<option value={jenis.id}>{jenis.nama || jenis.name}</option>
							{/each}
						</select>
						{#if validationErrors.jenis_situs_id}
							<span class="mt-1 text-xs text-error">{validationErrors.jenis_situs_id}</span>
						{/if}
					</div>
				</div>
			</div>

			<div class="card bg-base-100 shadow-md sm:shadow-xl">
				<div class="card-body p-4 sm:p-5">
					<h2 class="mb-3 card-title text-base sm:text-lg">Informasi Umum</h2>
					<div class="grid grid-cols-1 gap-3 sm:gap-4">
						<div class="form-control">
							<label class="label" for="nama">
								<span class="label-text font-medium"
									>Nama Situs <span class="text-error">*</span></span
								>
							</label>
							<input
								id="nama"
								type="text"
								class="input-bordered input min-h-11 w-full {validationErrors.nama
									? 'input-error'
									: ''}"
								bind:value={formData.nama}
								oninput={() => (validationErrors.nama = null)}
							/>
							{#if validationErrors.nama}
								<span class="mt-1 text-xs text-error">{validationErrors.nama}</span>
							{/if}
						</div>

						<div class="form-control">
							<label class="label" for="nomor_pengurus">
								<span class="label-text font-medium"
									>No. Telepon Pengurus <span class="text-error">*</span></span
								>
							</label>
							<div class="flex gap-2">
								<input
									id="nomor_pengurus"
									type="tel"
									placeholder="Contoh: 0812345678"
									class="input-bordered input min-h-11 flex-1 {validationErrors.nomor_telpon_pengurus
										? 'input-error'
										: ''}"
									bind:value={currentPhoneInput}
									onkeydown={(e) => {
										if (e.key === 'Enter') {
											e.preventDefault();
											addPhone();
										}
									}}
								/>
								<button type="button" class="btn btn-primary" onclick={addPhone}>Tambah</button>
							</div>

							<div class="mt-2 flex flex-wrap gap-2">
								{#each formData.nomor_telpon_pengurus as phone, index}
									<div class="badge gap-2 badge-lg badge-primary">
										{phone}
										<button
											type="button"
											onclick={() => removePhone(index)}
											class="btn btn-circle h-4 min-h-0 w-4 p-0 text-white btn-ghost btn-xs"
										>
											✕
										</button>
									</div>
								{/each}
							</div>

							{#if validationErrors.nomor_telpon_pengurus}
								<span class="mt-1 text-xs text-error">{validationErrors.nomor_telpon_pengurus}</span
								>
							{/if}
						</div>

						<div class="form-control">
							<label class="label" for="nomor_telepon">
								<span class="label-text font-medium"
									>Nomor Telepon Instansi/Situs <span class="text-xs text-base-content/50"
										>(Opsional)</span
									></span
								>
							</label>
							<input
								id="nomor_telepon"
								type="tel"
								class="input-bordered input min-h-11 w-full"
								bind:value={formData.nomor_telepon}
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
									value={formData.jenis_tipologi}
									readonly
								/>
							{:else}
								<select
									id="jenis_tipologi"
									class="select-bordered select min-h-11 w-full {validationErrors.jenis_tipologi
										? 'select-error'
										: ''}"
									bind:value={formData.jenis_tipologi}
									onchange={() => (validationErrors.jenis_tipologi = null)}
								>
									{#each getTipologiOptions() as option}
										<option value={option}>{option}</option>
									{/each}
								</select>
							{/if}
						</div>

						<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
							<div class="form-control">
								<label class="label" for="email">
									<span class="label-text font-medium"
										>Email <span class="text-xs text-base-content/50">(Opsional)</span></span
									>
								</label>
								<input
									id="email"
									type="email"
									class="input-bordered input min-h-11 w-full {validationErrors.email
										? 'input-error'
										: ''}"
									bind:value={formData.email}
									oninput={() => (validationErrors.email = null)}
								/>
								{#if validationErrors.email}<span class="mt-1 text-xs text-error"
										>{validationErrors.email}</span
									>{/if}
							</div>
							<div class="form-control">
								<label class="label" for="website">
									<span class="label-text font-medium"
										>Website <span class="text-xs text-base-content/50">(Opsional)</span></span
									>
								</label>
								<input
									id="website"
									type="url"
									class="input-bordered input min-h-11 w-full {validationErrors.website
										? 'input-error'
										: ''}"
									placeholder="https://..."
									bind:value={formData.website}
									oninput={() => (validationErrors.website = null)}
								/>
								{#if validationErrors.website}<span class="mt-1 text-xs text-error"
										>{validationErrors.website}</span
									>{/if}
							</div>
						</div>

						<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
							<div class="form-control">
								<label class="label" for="nomor_badan_hukum">
									<span class="label-text font-medium"
										>Nomor Badan Hukum <span class="text-xs text-base-content/50">(Opsional)</span
										></span
									>
								</label>
								<input
									id="nomor_badan_hukum"
									type="text"
									class="input-bordered input min-h-11 w-full"
									bind:value={formData.nomor_badan_hukum}
								/>
							</div>
							<div class="form-control">
								<label class="label" for="tahun_berdiri">
									<span class="label-text font-medium"
										>Tahun Berdiri <span class="text-error">*</span></span
									>
								</label>
								<input
									id="tahun_berdiri"
									type="number"
									class="input-bordered input min-h-11 w-full {validationErrors.tahun_berdiri
										? 'input-error'
										: ''}"
									bind:value={formData.tahun_berdiri}
									oninput={() => (validationErrors.tahun_berdiri = null)}
								/>
								{#if validationErrors.tahun_berdiri}<span class="mt-1 text-xs text-error"
										>{validationErrors.tahun_berdiri}</span
									>{/if}
							</div>
						</div>
					</div>
				</div>
			</div>

			<div class="card bg-base-100 shadow-md sm:shadow-xl">
				<div class="card-body p-4 sm:p-5">
					<h2 class="mb-3 card-title text-base sm:text-lg">Lokasi</h2>
					<div class="grid grid-cols-1 gap-3 sm:gap-4">
						<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
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
									id="desa"
									type="text"
									class="input-bordered input min-h-11 w-full {validationErrors.desa
										? 'input-error'
										: ''}"
									bind:value={formData.desa}
									oninput={() => (validationErrors.desa = null)}
								/>
								{#if validationErrors.desa}<span class="mt-1 text-xs text-error"
										>{validationErrors.desa}</span
									>{/if}
							</div>
						</div>

						<div class="form-control">
							<label class="label" for="alamat_lengkap">
								<span class="label-text font-medium"
									>Alamat Lengkap <span class="text-error">*</span></span
								>
							</label>
							<textarea
								id="alamat_lengkap"
								class="textarea-bordered textarea min-h-20 w-full {validationErrors.alamat_lengkap
									? 'textarea-error'
									: ''}"
								bind:value={formData.alamat_lengkap}
								oninput={() => (validationErrors.alamat_lengkap = null)}
								rows="2"
							></textarea>
							{#if validationErrors.alamat_lengkap}<span class="mt-1 text-xs text-error"
									>{validationErrors.alamat_lengkap}</span
								>{/if}
						</div>

						<div class="mt-4">
							<div class="mb-3 flex flex-col gap-2 sm:flex-row">
								<button
									type="button"
									class="btn min-h-11 flex-1 btn-primary sm:flex-none"
									onclick={getCurrentLocation}
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
									Lokasi Saat Ini
								</button>
							</div>

							<div class="mb-3 grid grid-cols-2 gap-2 sm:mb-4 sm:gap-4">
								<div class="form-control">
									<label class="label" for="latitude">
										<span class="label-text font-medium"
											>Latitude <span class="text-error">*</span></span
										>
									</label>
									<input
										required
										id="latitude"
										type="number"
										step="any"
										class="input-bordered input min-h-11 w-full {validationErrors.latitude
											? 'input-error'
											: ''}"
										value={formData.latitude ?? ''}
										readonly
									/>
								</div>
								<div class="form-control">
									<label class="label" for="longitude">
										<span class="label-text font-medium"
											>Longitude <span class="text-error">*</span></span
										>
									</label>
									<input
										required
										id="longitude"
										type="number"
										step="any"
										class="input-bordered input min-h-11 w-full {validationErrors.longitude
											? 'input-error'
											: ''}"
										value={formData.longitude ?? ''}
										readonly
									/>
								</div>
							</div>
							{#if validationErrors.latitude || validationErrors.longitude}
								<span class="mb-2 block text-xs text-error"
									>Silakan tentukan titik lokasi pada peta</span
								>
							{/if}

							<div
								class="h-64 w-full overflow-hidden rounded-lg border border-base-300 sm:h-80 md:h-96"
							>
								<div bind:this={mapContainer} class="h-full w-full"></div>
							</div>
						</div>
					</div>
				</div>
			</div>

			<div class="card bg-base-100 shadow-md sm:shadow-xl">
				<div class="card-body p-4 sm:p-5">
					<h2 class="mb-3 card-title text-base sm:text-lg">Kapasitas & Legalitas</h2>
					<div class="grid grid-cols-1 gap-3 sm:grid-cols-2 sm:gap-4">
						<div class="form-control">
							<label class="label" for="luas_tanah"
								><span class="label-text font-medium"
									>Luas Tanah (m²) <span class="text-error">*</span></span
								></label
							>
							<input
								id="luas_tanah"
								type="number"
								min="0"
								step="any"
								class="input-bordered input min-h-11 w-full {validationErrors.luas_tanah
									? 'input-error'
									: ''}"
								bind:value={formData.luas_tanah}
								oninput={() => (validationErrors.luas_tanah = null)}
							/>
							{#if validationErrors.luas_tanah}<span class="mt-1 text-xs text-error"
									>{validationErrors.luas_tanah}</span
								>{/if}
						</div>
						<div class="form-control">
							<label class="label" for="luas_bangunan"
								><span class="label-text font-medium"
									>Luas Bangunan (m²) <span class="text-error">*</span></span
								></label
							>
							<input
								id="luas_bangunan"
								type="number"
								min="0"
								step="any"
								class="input-bordered input min-h-11 w-full {validationErrors.luas_bangunan
									? 'input-error'
									: ''}"
								bind:value={formData.luas_bangunan}
								oninput={() => (validationErrors.luas_bangunan = null)}
							/>
							{#if validationErrors.luas_bangunan}<span class="mt-1 text-xs text-error"
									>{validationErrors.luas_bangunan}</span
								>{/if}
						</div>

						<div class="form-control">
							<label class="label" for="status_tanah"
								><span class="label-text font-medium"
									>Status Tanah <span class="text-error">*</span></span
								></label
							>
							<select
								id="status_tanah"
								class="select-bordered select min-h-11 w-full {validationErrors.status_tanah
									? 'select-error'
									: ''}"
								bind:value={formData.status_tanah}
								onchange={() => (validationErrors.status_tanah = null)}
							>
								<option value="">Pilih status...</option>
								<option value="Wakaf">Wakaf</option>
								<option value="Milik Sendiri">Milik Sendiri</option>
								<option value="Hibah">Hibah</option>
								<option value="Lainnya">Lainnya</option>
							</select>
							{#if validationErrors.status_tanah}<span class="mt-1 text-xs text-error"
									>{validationErrors.status_tanah}</span
								>{/if}
						</div>

						<div class="form-control">
							<label class="label" for="daya_tampung_max"
								><span class="label-text font-medium"
									>Daya Tampung Max <span class="text-error">*</span></span
								></label
							>
							<input
								id="daya_tampung_max"
								type="number"
								min="0"
								class="input-bordered input min-h-11 w-full {validationErrors.daya_tampung_max
									? 'input-error'
									: ''}"
								bind:value={formData.daya_tampung_max}
								oninput={() => (validationErrors.daya_tampung_max = null)}
							/>
							{#if validationErrors.daya_tampung_max}<span class="mt-1 text-xs text-error"
									>{validationErrors.daya_tampung_max}</span
								>{/if}
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
				<DetailMasjidForm bind:detail={formData.detail} bind:errors={validationErrors} />
			{/if}
			{#if isPesantren}
				<DetailPesantrenForm bind:detail={formData.detail} bind:errors={validationErrors} />
			{/if}
			{#if isMusholla}
				<DetailMusholaForm bind:detail={formData.detail} bind:errors={validationErrors} />
			{/if}
			{#if isMT}
				<DetailMTForm bind:detail={formData.detail} bind:errors={validationErrors} />
			{/if}

			<div class="card bg-base-100 shadow-md sm:shadow-xl">
				<div class="card-body p-4 sm:p-5">
					<h2 class="mb-3 card-title text-base sm:text-lg">Foto Situs</h2>
					<div class="form-control">
						<label class="label" for="images">
							<span class="label-text font-medium">Unggah Foto</span>
						</label>
						<div class="mb-2 flex flex-col gap-3">
							<input
								id="images"
								type="file"
								accept="image/*"
								multiple
								class="file-input-bordered file-input min-h-11 w-full"
								onchange={handleImageSelect}
								disabled={isConverting}
								bind:this={fileInputRef}
							/>
							{#if isConverting}
								<div class="flex items-center gap-2 text-sm text-base-content/70">
									<span class="loading loading-sm loading-spinner"></span> Mengkonversi gambar ke WebP...
								</div>
							{/if}
						</div>
						{#if selectedImages.length > 0}
							<div class="mt-2 grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4">
								{#each selectedImages as img, index (img.preview)}
									<div class="relative overflow-hidden rounded-lg border border-base-300">
										<img src={img.preview} alt="Preview" class="h-32 w-full object-cover" />
										<button
											type="button"
											class="btn absolute top-1 right-1 btn-circle btn-xs btn-error"
											onclick={() => removeImage(index)}>✕</button
										>
										<div
											class="absolute right-1 bottom-1 rounded bg-black/70 px-1 text-xs text-white"
										>
											WebP
										</div>
									</div>
								{/each}
							</div>
						{/if}
					</div>
				</div>
			</div>

			<div class="flex flex-col justify-end gap-3 pb-4 sm:flex-row sm:pb-8">
				<button
					onclick={() => goto('/dashboard/situs')}
					type="button"
					class="btn min-h-11 w-full btn-ghost sm:w-auto"
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
						Simpan Data
					{/if}
				</button>
			</div>
		</form>
	</div>

	<Toast
		show={showToast}
		message={toastMessage}
		type={toastType}
		onclose={() => (showToast = false)}
	/>
</div>

<style>
	:global(.leaflet-container) {
		z-index: 0;
	}
</style>
