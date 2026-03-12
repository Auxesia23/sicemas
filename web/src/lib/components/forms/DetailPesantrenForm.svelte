<script>
	// Svelte 5 Runes: Bindable prop for two-way binding with parent
	let { detail = $bindable() } = $props();

	// Local state for dynamic list inputs
	let newFasilitasPondok = $state('');

	// Initialize detail object if not exists
	function initDetail() {
		detail = {
			nama_yayasan: '',
			pimpinan_pesantren: '',
			kepengurusan: {
				ketua: '',
				sekretaris: '',
				bendahara: ''
			},
			santri: {
				mondok: {
					pria: null,
					wanita: null
				},
				tidak_mondok: {
					pria: null,
					wanita: null
				},
				total: {
					pria: null,
					wanita: null
				}
			},
			fasilitas_pondok: []
		};
	}

	// Ensure detail is initialized
	if (!detail?.nama_yayasan) {
		initDetail();
	}

	// Dynamic list helpers
	function addToArray(array, value, clearInput = true) {
		if (value.trim()) {
			array.push(value.trim());
			if (clearInput) return '';
		}
		return value;
	}

	function removeFromArray(array, index) {
		array.splice(index, 1);
	}

	function addFasilitasPondok() {
		newFasilitasPondok = addToArray(detail.fasilitas_pondok, newFasilitasPondok, true);
	}

	// Handle enter key for list inputs
	function handleEnterKey(event, addFn) {
		if (event.key === 'Enter') {
			event.preventDefault();
			addFn();
		}
	}

	// Helper functions for detail field bindings
	function updateKepengurusan(field, value) {
		if (!detail.kepengurusan) initDetail();
		detail.kepengurusan[field] = value;
	}

	function updateSantri(category, gender, value) {
		if (!detail.santri) initDetail();
		detail.santri[category][gender] = value;
	}

	function updateNamaYayasan(value) {
		detail.nama_yayasan = value;
	}

	function updatePimpinanPesantren(value) {
		detail.pimpinan_pesantren = value;
	}
</script>

<!-- Informasi Yayasan -->
<div class="card bg-base-100 shadow-md sm:shadow-xl">
	<div class="card-body p-4 sm:p-5">
		<h2 class="mb-3 card-title text-base sm:text-lg">Informasi Yayasan</h2>
		<div class="grid grid-cols-1 gap-3 sm:gap-4">
			<div class="form-control">
				<label class="label" for="nama_yayasan">
					<span class="label-text font-medium">Nama Yayasan</span>
				</label>
				<input
					id="nama_yayasan"
					type="text"
					class="input-bordered input min-h-11 w-full"
					placeholder="Yayasan Pendidikan Islam Darussalam"
					value={detail.nama_yayasan ?? ''}
					oninput={(e) => updateNamaYayasan(e.target.value)}
				/>
			</div>
			<div class="form-control">
				<label class="label" for="pimpinan_pesantren">
					<span class="label-text font-medium">Pimpinan Pesantren</span>
				</label>
				<input
					id="pimpinan_pesantren"
					type="text"
					class="input-bordered input min-h-11 w-full"
					placeholder="K.H. Ahmad Maulana"
					value={detail.pimpinan_pesantren ?? ''}
					oninput={(e) => updatePimpinanPesantren(e.target.value)}
				/>
			</div>
		</div>
	</div>
</div>

<!-- Kepengurusan -->
<div class="card bg-base-100 shadow-md sm:shadow-xl">
	<div class="card-body p-4 sm:p-5">
		<h2 class="mb-3 card-title text-base sm:text-lg">Kepengurusan</h2>
		<div class="grid grid-cols-1 gap-3 sm:grid-cols-3 sm:gap-4">
			<div class="form-control">
				<label class="label" for="ketua">
					<span class="label-text font-medium">Ketua</span>
				</label>
				<input
					id="ketua"
					type="text"
					class="input-bordered input min-h-11 w-full"
					placeholder="Ust. Zainuddin"
					value={detail.kepengurusan?.ketua ?? ''}
					oninput={(e) => updateKepengurusan('ketua', e.target.value)}
				/>
			</div>
			<div class="form-control">
				<label class="label" for="sekretaris">
					<span class="label-text font-medium">Sekretaris</span>
				</label>
				<input
					id="sekretaris"
					type="text"
					class="input-bordered input min-h-11 w-full"
					placeholder="Ust. Hasan Basri"
					value={detail.kepengurusan?.sekretaris ?? ''}
					oninput={(e) => updateKepengurusan('sekretaris', e.target.value)}
				/>
			</div>
			<div class="form-control">
				<label class="label" for="bendahara">
					<span class="label-text font-medium">Bendahara</span>
				</label>
				<input
					id="bendahara"
					type="text"
					class="input-bordered input min-h-11 w-full"
					placeholder="Hj. Fatimah"
					value={detail.kepengurusan?.bendahara ?? ''}
					oninput={(e) => updateKepengurusan('bendahara', e.target.value)}
				/>
			</div>
		</div>
	</div>
</div>

<!-- Data Santri -->
<div class="card bg-base-100 shadow-md sm:shadow-xl">
	<div class="card-body p-4 sm:p-5">
		<h2 class="mb-3 card-title text-base sm:text-lg">Data Santri</h2>
		<div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
			<!-- Mondok -->
			<div class="card bg-base-200">
				<div class="card-body p-3">
					<h3 class="card-title text-sm">Mondok</h3>
					<div class="mt-2 grid grid-cols-2 gap-2">
						<div class="form-control">
							<label class="label" for="mondok_pria">
								<span class="label-text text-xs">Pria</span>
							</label>
							<input
								id="mondok_pria"
								type="number"
								min="0"
								class="input-bordered input min-h-10 w-full"
								placeholder="0"
								value={detail.santri?.mondok?.pria ?? ''}
								oninput={(e) => updateSantri('mondok', 'pria', e.target.valueAsNumber || null)}
							/>
						</div>
						<div class="form-control">
							<label class="label" for="mondok_wanita">
								<span class="label-text text-xs">Wanita</span>
							</label>
							<input
								id="mondok_wanita"
								type="number"
								min="0"
								class="input-bordered input min-h-10 w-full"
								placeholder="0"
								value={detail.santri?.mondok?.wanita ?? ''}
								oninput={(e) => updateSantri('mondok', 'wanita', e.target.valueAsNumber || null)}
							/>
						</div>
					</div>
				</div>
			</div>

			<!-- Tidak Mondok -->
			<div class="card bg-base-200">
				<div class="card-body p-3">
					<h3 class="card-title text-sm">Tidak Mondok</h3>
					<div class="mt-2 grid grid-cols-2 gap-2">
						<div class="form-control">
							<label class="label" for="tidak_mondok_pria">
								<span class="label-text text-xs">Pria</span>
							</label>
							<input
								id="tidak_mondok_pria"
								type="number"
								min="0"
								class="input-bordered input min-h-10 w-full"
								placeholder="0"
								value={detail.santri?.tidak_mondok?.pria ?? ''}
								oninput={(e) => updateSantri('tidak_mondok', 'pria', e.target.valueAsNumber || null)}
							/>
						</div>
						<div class="form-control">
							<label class="label" for="tidak_mondok_wanita">
								<span class="label-text text-xs">Wanita</span>
							</label>
							<input
								id="tidak_mondok_wanita"
								type="number"
								min="0"
								class="input-bordered input min-h-10 w-full"
								placeholder="0"
								value={detail.santri?.tidak_mondok?.wanita ?? ''}
								oninput={(e) => updateSantri('tidak_mondok', 'wanita', e.target.valueAsNumber || null)}
							/>
						</div>
					</div>
				</div>
			</div>

			<!-- Total -->
			<div class="card bg-base-200">
				<div class="card-body p-3">
					<h3 class="card-title text-sm">Total</h3>
					<div class="mt-2 grid grid-cols-2 gap-2">
						<div class="form-control">
							<label class="label" for="total_pria">
								<span class="label-text text-xs">Pria</span>
							</label>
							<input
								id="total_pria"
								type="number"
								min="0"
								class="input-bordered input min-h-10 w-full"
								placeholder="0"
								value={detail.santri?.total?.pria ?? ''}
								oninput={(e) => updateSantri('total', 'pria', e.target.valueAsNumber || null)}
							/>
						</div>
						<div class="form-control">
							<label class="label" for="total_wanita">
								<span class="label-text text-xs">Wanita</span>
							</label>
							<input
								id="total_wanita"
								type="number"
								min="0"
								class="input-bordered input min-h-10 w-full"
								placeholder="0"
								value={detail.santri?.total?.wanita ?? ''}
								oninput={(e) => updateSantri('total', 'wanita', e.target.valueAsNumber || null)}
							/>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>

<!-- Fasilitas Pondok -->
<div class="card bg-base-100 shadow-md sm:shadow-xl">
	<div class="card-body p-4 sm:p-5">
		<h2 class="mb-3 card-title text-base sm:text-lg">Fasilitas Pondok</h2>
		<div class="mb-3 sm:mb-4">
			<label class="label" for="fasilitas_pondok_input">
				<span class="label-text font-medium">Tambah Fasilitas</span>
			</label>
			<div class="mb-2 flex flex-col gap-2 sm:flex-row">
				<input
					id="fasilitas_pondok_input"
					type="text"
					class="input-bordered input min-h-11 flex-1"
					placeholder="Tambah fasilitas..."
					bind:value={newFasilitasPondok}
					onkeydown={(e) => handleEnterKey(e, addFasilitasPondok)}
				/>
				<button type="button" class="btn min-h-11 btn-secondary" onclick={addFasilitasPondok}>
					Tambah
				</button>
			</div>
			<div class="flex flex-wrap gap-2">
				{#each detail.fasilitas_pondok || [] as fasilitas, index (index)}
					<div class="badge min-h-8 gap-2 badge-primary">
						{fasilitas}
						<button
							type="button"
							class="btn btn-ghost btn-xs"
							onclick={() => removeFromArray(detail.fasilitas_pondok, index)}
						>
							✕
						</button>
					</div>
				{/each}
				{#if detail.fasilitas_pondok?.length === 0}
					<span class="text-sm text-base-content/50">Belum ada fasilitas pondok</span>
				{/if}
			</div>
		</div>
	</div>
</div>
