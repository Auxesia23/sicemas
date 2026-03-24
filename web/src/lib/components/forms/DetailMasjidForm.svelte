<script>
	// Tambahkan prop errors untuk menangkap pesan validasi dari parent
	// Gunakan $bindable agar sinkronisasi data lancar
	let { detail = $bindable(), errors = $bindable({}) } = $props();

	// Local state for dynamic list inputs (Badge System)
	let newJenisBuku = $state('');
	let newNamaImam = $state('');
	let newNamaMuazdin = $state('');
	let newSekretaris = $state('');
	let newBendahara = $state('');
	let newFasilitasUmum = $state('');
	let newFasilitasRamanAnak = $state('');
	let newFasilitasDisabilitas = $state('');
	let newKegiatan = $state('');

	function initDetail() {
		detail = {
			perpustakaan: {
				luas_m2: 0,
				jumlah_pengurus: 0,
				kondisi: '',
				jenis_buku: []
			},
			kalibrasi_arah_kiblat: {
				azimut: '',
				tanggal_kalibrasi: ''
			},
			sdm_masjid: {
				jumlah_pengurus: 0,
				jumlah_imam: 0,
				jumlah_khotib: 0,
				jumlah_muadzin: 0,
				jumlah_remaja: 0,
				jumlah_jemaah: 0
			},
			pengurus_dkm: {
				ketua: '',
				sekretaris: [],
				bendahara: [],
				nama_imam: [],
				nama_muazdin: []
			},
			fasilitas_umum: [],
			fasilitas_raman_anak: [],
			fasilitas_disabilitas: [],
			kegiatan: []
		};
	}

	if (!detail?.perpustakaan) {
		initDetail();
	}

	// Fix data lama jika sekretaris/bendahara masih berupa string
	if (typeof detail.pengurus_dkm?.sekretaris === 'string')
		detail.pengurus_dkm.sekretaris = [detail.pengurus_dkm.sekretaris].filter(Boolean);
	if (typeof detail.pengurus_dkm?.bendahara === 'string')
		detail.pengurus_dkm.bendahara = [detail.pengurus_dkm.bendahara].filter(Boolean);

	// Helper penambah array (Badge)
	function addToArray(array, value) {
		if (value.trim() && !array.includes(value.trim())) {
			array.push(value.trim());
			return ''; // Reset input field
		}
		return value;
	}

	function removeFromArray(array, index) {
		array.splice(index, 1);
	}

	// Event Handlers
	function handleEnterKey(event, addFn) {
		if (event.key === 'Enter') {
			event.preventDefault();
			addFn();
		}
	}

	// Fungsi pembantu untuk hapus error saat mengetik
	function clearError(parent, child) {
		if (errors?.[parent]?.[child]) {
			errors[parent][child] = null;
		}
	}
</script>

<div class="card mt-4 bg-base-100 shadow-md sm:shadow-xl">
	<div class="card-body p-4 sm:p-5">
		<h2 class="mb-3 card-title text-base text-primary sm:text-lg">
			SDM Masjid <span class="text-error">*</span>
		</h2>
		{#if errors?.sdm_masjid && typeof errors.sdm_masjid === 'string'}
			<div class="mb-2 text-sm text-error">{errors.sdm_masjid}</div>
		{/if}
		<div class="grid grid-cols-2 gap-2 sm:grid-cols-3 sm:gap-4">
			{#each [{ id: 'pengurus', label: 'Jml Pengurus', path: 'jumlah_pengurus' }, { id: 'imam', label: 'Jml Imam', path: 'jumlah_imam' }, { id: 'khotib', label: 'Jml Khotib', path: 'jumlah_khotib' }, { id: 'muadzin', label: 'Jml Muadzin', path: 'jumlah_muadzin' }, { id: 'remaja', label: 'Jml Remaja', path: 'jumlah_remaja' }, { id: 'jemaah', label: 'Jml Jemaah', path: 'jumlah_jemaah' }] as item}
				<div class="form-control">
					<label class="label" for="sdm_{item.id}">
						<span class="label-text font-medium"
							>{item.label} <span class="text-error">*</span></span
						>
					</label>
					<input
						id="sdm_{item.id}"
						type="number"
						min="0"
						class="input-bordered input min-h-11 w-full {errors?.sdm_masjid?.[item.path]
							? 'input-error'
							: ''}"
						bind:value={detail.sdm_masjid[item.path]}
						oninput={() => clearError('sdm_masjid', item.path)}
					/>
					{#if errors?.sdm_masjid?.[item.path]}
						<span class="mt-1 text-[10px] text-error">{errors.sdm_masjid[item.path]}</span>
					{/if}
				</div>
			{/each}
		</div>
	</div>
</div>

<div class="card mt-4 bg-base-100 shadow-md sm:shadow-xl">
	<div class="card-body p-4 sm:p-5">
		<h2 class="mb-3 card-title text-base text-primary sm:text-lg">
			Pengurus DKM <span class="text-error">*</span>
		</h2>

		<div class="form-control mb-4">
			<label class="label" for="pengurus_ketua">
				<span class="label-text font-medium">Nama Ketua DKM <span class="text-error">*</span></span>
			</label>
			<input
				id="pengurus_ketua"
				type="text"
				class="input-bordered input min-h-11 w-full {errors?.pengurus_dkm?.ketua
					? 'input-error'
					: ''}"
				bind:value={detail.pengurus_dkm.ketua}
				oninput={() => clearError('pengurus_dkm', 'ketua')}
			/>
			{#if errors?.pengurus_dkm?.ketua}<span class="mt-1 text-xs text-error"
					>{errors.pengurus_dkm.ketua}</span
				>{/if}
		</div>

		<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
			<div class="form-control">
				<label class="label" for="sek_in"
					><span class="label-text font-medium">Sekretaris <span class="text-error">*</span></span
					></label
				>
				<div class="flex gap-2">
					<input
						id="sek_in"
						type="text"
						class="input-bordered input flex-1 {errors?.pengurus_dkm?.sekretaris
							? 'input-error'
							: ''}"
						placeholder="Tambah Sekretaris..."
						bind:value={newSekretaris}
						onkeydown={(e) =>
							handleEnterKey(
								e,
								() => (newSekretaris = addToArray(detail.pengurus_dkm.sekretaris, newSekretaris))
							)}
					/>
					<button
						type="button"
						class="btn btn-secondary"
						onclick={() =>
							(newSekretaris = addToArray(detail.pengurus_dkm.sekretaris, newSekretaris))}
						>Tambah</button
					>
				</div>
				<div class="mt-2 flex flex-wrap gap-2">
					{#each detail.pengurus_dkm.sekretaris as p, i}
						<div class="badge h-8 gap-2 text-white badge-info">
							{p}<button
								type="button"
								onclick={() => removeFromArray(detail.pengurus_dkm.sekretaris, i)}>✕</button
							>
						</div>
					{/each}
				</div>
				{#if errors?.pengurus_dkm?.sekretaris}<span class="mt-1 text-xs text-error"
						>{errors.pengurus_dkm.sekretaris}</span
					>{/if}
			</div>

			<div class="form-control">
				<label class="label" for="ben_in"
					><span class="label-text font-medium">Bendahara <span class="text-error">*</span></span
					></label
				>
				<div class="flex gap-2">
					<input
						id="ben_in"
						type="text"
						class="input-bordered input flex-1 {errors?.pengurus_dkm?.bendahara
							? 'input-error'
							: ''}"
						placeholder="Tambah Bendahara..."
						bind:value={newBendahara}
						onkeydown={(e) =>
							handleEnterKey(
								e,
								() => (newBendahara = addToArray(detail.pengurus_dkm.bendahara, newBendahara))
							)}
					/>
					<button
						type="button"
						class="btn btn-secondary"
						onclick={() => (newBendahara = addToArray(detail.pengurus_dkm.bendahara, newBendahara))}
						>Tambah</button
					>
				</div>
				<div class="mt-2 flex flex-wrap gap-2">
					{#each detail.pengurus_dkm.bendahara as p, i}
						<div class="badge h-8 gap-2 badge-warning">
							{p}<button
								type="button"
								onclick={() => removeFromArray(detail.pengurus_dkm.bendahara, i)}>✕</button
							>
						</div>
					{/each}
				</div>
				{#if errors?.pengurus_dkm?.bendahara}<span class="mt-1 text-xs text-error"
						>{errors.pengurus_dkm.bendahara}</span
					>{/if}
			</div>
		</div>

		<div class="mt-4 grid grid-cols-1 gap-4 sm:grid-cols-2">
			<div class="form-control">
				<label class="label" for="imam_in"
					><span class="label-text font-medium">Nama Imam <span class="text-error">*</span></span
					></label
				>
				<div class="flex gap-2">
					<input
						id="imam_in"
						type="text"
						class="input-bordered input flex-1 {errors?.pengurus_dkm?.nama_imam
							? 'input-error'
							: ''}"
						bind:value={newNamaImam}
						onkeydown={(e) =>
							handleEnterKey(
								e,
								() => (newNamaImam = addToArray(detail.pengurus_dkm.nama_imam, newNamaImam))
							)}
					/>
					<button
						type="button"
						class="btn btn-secondary"
						onclick={() => (newNamaImam = addToArray(detail.pengurus_dkm.nama_imam, newNamaImam))}
						>Tambah</button
					>
				</div>
				<div class="mt-2 flex flex-wrap gap-2">
					{#each detail.pengurus_dkm.nama_imam as p, i}
						<div class="badge h-8 gap-2 badge-primary">
							{p}<button
								type="button"
								onclick={() => removeFromArray(detail.pengurus_dkm.nama_imam, i)}>✕</button
							>
						</div>
					{/each}
				</div>
				{#if errors?.pengurus_dkm?.nama_imam}<span class="mt-1 text-xs text-error"
						>{errors.pengurus_dkm.nama_imam}</span
					>{/if}
			</div>

			<div class="form-control">
				<label class="label" for="mua_in"
					><span class="label-text font-medium">Nama Muadzin <span class="text-error">*</span></span
					></label
				>
				<div class="flex gap-2">
					<input
						id="mua_in"
						type="text"
						class="input-bordered input flex-1 {errors?.pengurus_dkm?.nama_muazdin
							? 'input-error'
							: ''}"
						bind:value={newNamaMuazdin}
						onkeydown={(e) =>
							handleEnterKey(
								e,
								() =>
									(newNamaMuazdin = addToArray(detail.pengurus_dkm.nama_muazdin, newNamaMuazdin))
							)}
					/>
					<button
						type="button"
						class="btn btn-secondary"
						onclick={() =>
							(newNamaMuazdin = addToArray(detail.pengurus_dkm.nama_muazdin, newNamaMuazdin))}
						>Tambah</button
					>
				</div>
				<div class="mt-2 flex flex-wrap gap-2">
					{#each detail.pengurus_dkm.nama_muazdin as p, i}
						<div class="badge h-8 gap-2 badge-primary">
							{p}<button
								type="button"
								onclick={() => removeFromArray(detail.pengurus_dkm.nama_muazdin, i)}>✕</button
							>
						</div>
					{/each}
				</div>
				{#if errors?.pengurus_dkm?.nama_muazdin}<span class="mt-1 text-xs text-error"
						>{errors.pengurus_dkm.nama_muazdin}</span
					>{/if}
			</div>
		</div>
	</div>
</div>

<div class="card mt-4 bg-base-100 shadow-md sm:shadow-xl">
	<div class="card-body p-4 sm:p-5">
		<h2 class="mb-3 card-title text-base sm:text-lg">Perpustakaan</h2>
		<div class="grid grid-cols-2 gap-2 sm:grid-cols-3 sm:gap-4">
			<div class="form-control">
				<label class="label" for="perpustakaan_luas_m2"
					><span class="label-text font-medium">Luas (m²)</span></label
				>
				<input
					id="perpustakaan_luas_m2"
					type="number"
					min="0"
					class="input-bordered input w-full"
					bind:value={detail.perpustakaan.luas_m2}
				/>
			</div>
			<div class="form-control">
				<label class="label" for="perpustakaan_jumlah_pengurus"
					><span class="label-text font-medium">Jml. Pengurus</span></label
				>
				<input
					id="perpustakaan_jumlah_pengurus"
					type="number"
					min="0"
					class="input-bordered input w-full"
					bind:value={detail.perpustakaan.jumlah_pengurus}
				/>
			</div>
			<div class="form-control col-span-2 sm:col-span-1">
				<label class="label" for="perpustakaan_kondisi"
					><span class="label-text font-medium">Kondisi</span></label
				>
				<select
					id="perpustakaan_kondisi"
					class="select-bordered select w-full"
					bind:value={detail.perpustakaan.kondisi}
				>
					<option value="">Pilih...</option>
					<option value="Baik">Baik</option>
					<option value="Cukup Baik">Cukup Baik</option>
					<option value="Kurang Baik">Kurang Baik</option>
					<option value="Rusak">Rusak</option>
				</select>
			</div>
		</div>

		<div class="mt-4">
			<label class="label" for="jenis_buku_input"
				><span class="label-text font-medium">Jenis Buku</span></label
			>
			<div class="mb-2 flex gap-2">
				<input
					id="jenis_buku_input"
					type="text"
					class="input-bordered input flex-1"
					placeholder="Tambah jenis buku..."
					bind:value={newJenisBuku}
					onkeydown={(e) =>
						handleEnterKey(
							e,
							() => (newJenisBuku = addToArray(detail.perpustakaan.jenis_buku, newJenisBuku))
						)}
				/>
				<button
					type="button"
					class="btn btn-secondary"
					onclick={() => (newJenisBuku = addToArray(detail.perpustakaan.jenis_buku, newJenisBuku))}
					>Tambah</button
				>
			</div>
			<div class="flex flex-wrap gap-2">
				{#each detail.perpustakaan.jenis_buku as buku, index}
					<div class="badge h-8 gap-2 badge-primary">
						{buku}<button
							type="button"
							onclick={() => removeFromArray(detail.perpustakaan.jenis_buku, index)}>✕</button
						>
					</div>
				{/each}
			</div>
		</div>
	</div>
</div>

<div class="card mt-4 bg-base-100 shadow-md sm:shadow-xl">
	<div class="card-body p-4 sm:p-5">
		<h2 class="mb-3 card-title text-base sm:text-lg">Kalibrasi Kiblat</h2>
		<div class="grid grid-cols-1 gap-3 sm:grid-cols-2 sm:gap-4">
			<div class="form-control">
				<label class="label" for="kalibrasi_azimut"
					><span class="label-text font-medium">Azimut</span></label
				>
				<input
					id="kalibrasi_azimut"
					type="text"
					class="input-bordered input w-full"
					placeholder="295.12"
					bind:value={detail.kalibrasi_arah_kiblat.azimut}
				/>
			</div>
			<div class="form-control">
				<label class="label" for="kalibrasi_tanggal"
					><span class="label-text font-medium">Tanggal Kalibrasi</span></label
				>
				<input
					id="kalibrasi_tanggal"
					type="date"
					class="input-bordered input w-full"
					bind:value={detail.kalibrasi_arah_kiblat.tanggal_kalibrasi}
				/>
			</div>
		</div>
	</div>
</div>

<div class="card mt-4 bg-base-100 shadow-md sm:shadow-xl">
	<div class="card-body p-4 sm:p-5">
		<h2 class="mb-3 card-title text-base sm:text-lg">Fasilitas & Kegiatan</h2>

		<div class="mb-4">
			<label class="label" for="fu_in"
				><span class="label-text font-medium">Fasilitas Umum</span></label
			>
			<div class="mb-2 flex gap-2">
				<input
					id="fu_in"
					type="text"
					class="input-bordered input flex-1"
					placeholder="Contoh: Area Parkir"
					bind:value={newFasilitasUmum}
					onkeydown={(e) =>
						handleEnterKey(
							e,
							() => (newFasilitasUmum = addToArray(detail.fasilitas_umum, newFasilitasUmum))
						)}
				/>
				<button
					type="button"
					class="btn btn-secondary"
					onclick={() => (newFasilitasUmum = addToArray(detail.fasilitas_umum, newFasilitasUmum))}
					>Tambah</button
				>
			</div>
			<div class="flex flex-wrap gap-2">
				{#each detail.fasilitas_umum as f, i}
					<div class="badge h-8 gap-2 badge-secondary">
						{f}<button type="button" onclick={() => removeFromArray(detail.fasilitas_umum, i)}
							>✕</button
						>
					</div>
				{/each}
			</div>
		</div>

		<div class="mb-4">
			<label class="label" for="fra_in"
				><span class="label-text font-medium">Fasilitas Ramah Anak</span></label
			>
			<div class="mb-2 flex gap-2">
				<input
					id="fra_in"
					type="text"
					class="input-bordered input flex-1"
					placeholder="Contoh: Ruang Bermain"
					bind:value={newFasilitasRamanAnak}
					onkeydown={(e) =>
						handleEnterKey(
							e,
							() =>
								(newFasilitasRamanAnak = addToArray(
									detail.fasilitas_raman_anak,
									newFasilitasRamanAnak
								))
						)}
				/>
				<button
					type="button"
					class="btn btn-secondary"
					onclick={() =>
						(newFasilitasRamanAnak = addToArray(
							detail.fasilitas_raman_anak,
							newFasilitasRamanAnak
						))}>Tambah</button
				>
			</div>
			<div class="flex flex-wrap gap-2">
				{#each detail.fasilitas_raman_anak as f, i}
					<div class="badge h-8 gap-2 badge-accent">
						{f}<button type="button" onclick={() => removeFromArray(detail.fasilitas_raman_anak, i)}
							>✕</button
						>
					</div>
				{/each}
			</div>
		</div>

		<div>
			<label class="label" for="keg_in"
				><span class="label-text font-medium">Kegiatan Rutin</span></label
			>
			<div class="mb-2 flex gap-2">
				<input
					id="keg_in"
					type="text"
					class="input-bordered input flex-1"
					placeholder="Contoh: Pengajian Mingguan"
					bind:value={newKegiatan}
					onkeydown={(e) =>
						handleEnterKey(e, () => (newKegiatan = addToArray(detail.kegiatan, newKegiatan)))}
				/>
				<button
					type="button"
					class="btn btn-secondary"
					onclick={() => (newKegiatan = addToArray(detail.kegiatan, newKegiatan))}>Tambah</button
				>
			</div>
			<div class="flex flex-wrap gap-2">
				{#each detail.kegiatan as k, i}
					<div class="badge h-8 gap-2 badge-success">
						{k}<button type="button" onclick={() => removeFromArray(detail.kegiatan, i)}>✕</button>
					</div>
				{/each}
			</div>
		</div>
	</div>
</div>
