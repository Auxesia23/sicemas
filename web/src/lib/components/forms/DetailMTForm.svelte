<script>
	// Svelte 5 Runes: Bindable prop for two-way binding with parent
	let { detail = $bindable() } = $props();

	// Local state for penceramah inputs
	let newPenceramahNama = $state('');
	let newPenceramahMateri = $state('');

	// Initialize detail object if not exists
	function initDetail() {
		detail = {
			pengurus: {
				ketua: '',
				sekretaris: '',
				bendahara: ''
			},
			kehadiran_pria: {
				jumlah_dewasa: null,
				jumlah_remaja: null,
				waktu_pengajian: ''
			},
			kehadiran_wanita: {
				jumlah_dewasa: null,
				jumlah_remaja: null,
				waktu_pengajian: ''
			},
			penceramah: []
		};
	}

	// Ensure detail is initialized (synchronous, NOT in $effect)
	if (!detail || !detail.pengurus) {
		initDetail();
	}

	// Helper functions for detail field bindings
	function updatePengurus(field, value) {
		if (!detail.pengurus) initDetail();
		detail.pengurus[field] = value;
	}

	function updateKehadiranPria(field, value) {
		if (!detail.kehadiran_pria) initDetail();
		detail.kehadiran_pria[field] = value;
	}

	function updateKehadiranWanita(field, value) {
		if (!detail.kehadiran_wanita) initDetail();
		detail.kehadiran_wanita[field] = value;
	}

	// Add penceramah to the list
	function addPenceramah() {
		if (newPenceramahNama.trim() && newPenceramahMateri.trim()) {
			detail.penceramah.push({
				nama: newPenceramahNama.trim(),
				materi: newPenceramahMateri.trim()
			});
			newPenceramahNama = '';
			newPenceramahMateri = '';
		}
	}

	// Remove penceramah from the list
	function removePenceramah(index) {
		detail.penceramah.splice(index, 1);
	}

	// Handle enter key for inputs
	function handleEnterKey(event) {
		if (event.key === 'Enter') {
			event.preventDefault();
			addPenceramah();
		}
	}
</script>

<!-- Pengurus -->
<div class="card bg-base-100 shadow-md sm:shadow-xl">
	<div class="card-body p-4 sm:p-5">
		<h2 class="mb-3 card-title text-base sm:text-lg">Pengurus</h2>
		<div class="grid grid-cols-1 gap-3 sm:grid-cols-3 sm:gap-4">
			<div class="form-control">
				<label class="label" for="pengurus_ketua">
					<span class="label-text font-medium">Ketua</span>
				</label>
				<input
					id="pengurus_ketua"
					type="text"
					class="input-bordered input min-h-11 w-full"
					placeholder="Hj. Siti Aminah"
					value={detail.pengurus?.ketua ?? ''}
					oninput={(e) => updatePengurus('ketua', e.target.value)}
				/>
			</div>
			<div class="form-control">
				<label class="label" for="pengurus_sekretaris">
					<span class="label-text font-medium">Sekretaris</span>
				</label>
				<input
					id="pengurus_sekretaris"
					type="text"
					class="input-bordered input min-h-11 w-full"
					placeholder="Ibu Nurul Aini"
					value={detail.pengurus?.sekretaris ?? ''}
					oninput={(e) => updatePengurus('sekretaris', e.target.value)}
				/>
			</div>
			<div class="form-control">
				<label class="label" for="pengurus_bendahara">
					<span class="label-text font-medium">Bendahara</span>
				</label>
				<input
					id="pengurus_bendahara"
					type="text"
					class="input-bordered input min-h-11 w-full"
					placeholder="Ibu Dewi Sartika"
					value={detail.pengurus?.bendahara ?? ''}
					oninput={(e) => updatePengurus('bendahara', e.target.value)}
				/>
			</div>
		</div>
	</div>
</div>

<!-- Kehadiran Pria -->
<div class="card bg-base-100 shadow-md sm:shadow-xl">
	<div class="card-body p-4 sm:p-5">
		<h2 class="mb-3 card-title text-base sm:text-lg">Kehadiran Pria</h2>
		<div class="grid grid-cols-1 gap-3 sm:grid-cols-3 sm:gap-4">
			<div class="form-control">
				<label class="label" for="pria_dewasa">
					<span class="label-text font-medium">Jml. Dewasa</span>
				</label>
				<input
					id="pria_dewasa"
					type="number"
					min="0"
					class="input-bordered input min-h-11 w-full"
					placeholder="0"
					value={detail.kehadiran_pria?.jumlah_dewasa ?? ''}
					oninput={(e) => updateKehadiranPria('jumlah_dewasa', e.target.valueAsNumber || null)}
				/>
			</div>
			<div class="form-control">
				<label class="label" for="pria_remaja">
					<span class="label-text font-medium">Jml. Remaja</span>
				</label>
				<input
					id="pria_remaja"
					type="number"
					min="0"
					class="input-bordered input min-h-11 w-full"
					placeholder="0"
					value={detail.kehadiran_pria?.jumlah_remaja ?? ''}
					oninput={(e) => updateKehadiranPria('jumlah_remaja', e.target.valueAsNumber || null)}
				/>
			</div>
			<div class="form-control">
				<label class="label" for="pria_waktu">
					<span class="label-text font-medium">Waktu Pengajian</span>
				</label>
				<select
					id="pria_waktu"
					class="select-bordered select min-h-11 w-full"
					value={detail.kehadiran_pria?.waktu_pengajian ?? ''}
					oninput={(e) => updateKehadiranPria('waktu_pengajian', e.target.value)}
				>
					<option value="">Pilih...</option>
					<option value="Pagi">Pagi</option>
					<option value="Sore">Sore</option>
					<option value="Malam">Malam</option>
				</select>
			</div>
		</div>
	</div>
</div>

<!-- Kehadiran Wanita -->
<div class="card bg-base-100 shadow-md sm:shadow-xl">
	<div class="card-body p-4 sm:p-5">
		<h2 class="mb-3 card-title text-base sm:text-lg">Kehadiran Wanita</h2>
		<div class="grid grid-cols-1 gap-3 sm:grid-cols-3 sm:gap-4">
			<div class="form-control">
				<label class="label" for="wanita_dewasa">
					<span class="label-text font-medium">Jml. Dewasa</span>
				</label>
				<input
					id="wanita_dewasa"
					type="number"
					min="0"
					class="input-bordered input min-h-11 w-full"
					placeholder="0"
					value={detail.kehadiran_wanita?.jumlah_dewasa ?? ''}
					oninput={(e) => updateKehadiranWanita('jumlah_dewasa', e.target.valueAsNumber || null)}
				/>
			</div>
			<div class="form-control">
				<label class="label" for="wanita_remaja">
					<span class="label-text font-medium">Jml. Remaja</span>
				</label>
				<input
					id="wanita_remaja"
					type="number"
					min="0"
					class="input-bordered input min-h-11 w-full"
					placeholder="0"
					value={detail.kehadiran_wanita?.jumlah_remaja ?? ''}
					oninput={(e) => updateKehadiranWanita('jumlah_remaja', e.target.valueAsNumber || null)}
				/>
			</div>
			<div class="form-control">
				<label class="label" for="wanita_waktu">
					<span class="label-text font-medium">Waktu Pengajian</span>
				</label>
				<select
					id="wanita_waktu"
					class="select-bordered select min-h-11 w-full"
					value={detail.kehadiran_wanita?.waktu_pengajian ?? ''}
					oninput={(e) => updateKehadiranWanita('waktu_pengajian', e.target.value)}
				>
					<option value="">Pilih...</option>
					<option value="Pagi">Pagi</option>
					<option value="Sore">Sore</option>
				</select>
			</div>
		</div>
	</div>
</div>

<!-- Daftar Penceramah -->
<div class="card bg-base-100 shadow-md sm:shadow-xl">
	<div class="card-body p-4 sm:p-5">
		<h2 class="mb-3 card-title text-base sm:text-lg">Daftar Penceramah</h2>
		<div class="mb-3 sm:mb-4">
			<label class="label" for="penceramah_inputs">
				<span class="label-text font-medium">Tambah Penceramah</span>
			</label>
			<div class="mb-2 grid grid-cols-1 gap-2 sm:grid-cols-12">
				<div class="form-control sm:col-span-5">
					<input
						id="penceramah_nama"
						type="text"
						class="input-bordered input min-h-11 w-full"
						placeholder="Nama penceramah..."
						bind:value={newPenceramahNama}
						onkeydown={handleEnterKey}
					/>
				</div>
				<div class="form-control sm:col-span-5">
					<input
						id="penceramah_materi"
						type="text"
						class="input-bordered input min-h-11 w-full"
						placeholder="Materi..."
						bind:value={newPenceramahMateri}
						onkeydown={handleEnterKey}
					/>
				</div>
				<button type="button" class="btn min-h-11 btn-secondary sm:col-span-2" onclick={addPenceramah}>
					Tambah
				</button>
			</div>

			<!-- Penceramah List -->
			<div class="flex flex-col gap-2">
				{#each detail.penceramah || [] as penceramah, index (index)}
					<div class="flex items-center justify-between rounded-lg border border-base-300 bg-base-200 p-3">
						<div class="flex-1">
							<div class="font-medium">{penceramah.nama}</div>
							<div class="text-sm text-base-content/70">{penceramah.materi}</div>
						</div>
						<button
							type="button"
							class="btn btn-circle btn-xs btn-error"
							onclick={() => removePenceramah(index)}
						>
							✕
						</button>
					</div>
				{/each}
				{#if detail.penceramah?.length === 0}
					<span class="text-sm text-base-content/50">Belum ada penceramah</span>
				{/if}
			</div>
		</div>
	</div>
</div>
