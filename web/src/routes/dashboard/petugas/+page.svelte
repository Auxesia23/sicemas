<script>
	import { onMount } from 'svelte';
	import { z } from 'zod';
	import apiService from '$lib/api';
	import { hasAllPermissions } from '$lib/permissions.js';

	let { data } = $props();
	let user = $derived(data.user);

	let staffs = $state([]);
	let isLoading = $state(false);
	let isModalOpen = $state(false);
	let isEditMode = $state(false);
	let editingUserId = $state(null);
	let roles = $state([]);

	// State untuk pesan error/sukses
	let tableError = $state('');
	let modalError = $state('');
	let successMessage = $state('');

	// Form state dengan Jabatan & Unit Kerja (Opsional)
	let formData = $state({
		nip: '',
		nama_lengkap: '',
		email: '',
		nomor_telepon: '',
		jabatan: '',
		unit_kerja: '',
		peran: ''
	});

	// Zod validation schema
	const userSchema = z.object({
		nip: z.string().length(18, 'NIP harus berjumlah 18 digit'),
		nama_lengkap: z.string().min(1, 'Nama lengkap wajib diisi'),
		email: z.string().email('Format email tidak valid'),
		nomor_telepon: z.string().min(1, 'Nomor telepon wajib diisi'),
		jabatan: z.string().min(1, 'Jabatan wajib diisi'),
		unit_kerja: z.string().optional(),
		peran: z.string().min(1, 'Peran wajib dipilih')
	});

	onMount(async () => {
		await getAllUsers();
		await getRoles();
	});

	const getAllUsers = async () => {
		isLoading = true;
		tableError = '';
		try {
			const response = await apiService.get('/users');
			if (!response.ok) {
				const textError = await response.text();
				throw new Error(textError || 'Gagal mengambil data petugas');
			}
			staffs = await response.json();
		} catch (error) {
			tableError = error.message;
		} finally {
			isLoading = false;
		}
	};

	const getRoles = async () => {
		try {
			const response = await apiService.get('/roles');
			if (!response.ok) {
				throw new Error('Gagal mengambil data peran');
			}
			roles = await response.json();
		} catch (error) {
			console.error('Failed to fetch roles:', error);
		}
	};

	const openModal = () => {
		isModalOpen = true;
		isEditMode = false;
		editingUserId = null;
		modalError = '';
		successMessage = '';
		formData = {
			nip: '',
			nama_lengkap: '',
			email: '',
			nomor_telepon: '',
			jabatan: '',
			unit_kerja: '',
			peran: ''
		};
	};

	const openEditModal = (user) => {
		isModalOpen = true;
		isEditMode = true;
		editingUserId = user.id;
		modalError = '';
		successMessage = '';
		formData = {
			nip: user.nip || '',
			nama_lengkap: user.nama_lengkap || '',
			email: user.email || '',
			nomor_telepon: user.nomor_telepon || '',
			jabatan: user.jabatan || '',
			unit_kerja: user.unit_kerja || '',
			peran: user.peran || ''
		};
	};

	const closeModal = () => {
		isModalOpen = false;
		isEditMode = false;
		editingUserId = null;
		modalError = '';
	};

	const handleSubmit = async () => {
		modalError = '';
		successMessage = '';

		// 1. Validasi Client-side dengan Zod
		const validation = userSchema.safeParse(formData);
		if (!validation.success) {
			modalError = validation.error.issues[0].message;
			return;
		}

		isLoading = true;
		try {
			let response;
			if (isEditMode) {
				// Edit mode - PUT request
				response = await apiService.put(`/users/${editingUserId}`, formData);
			} else {
				// Add mode - POST request
				response = await apiService.post('/users', formData);
			}

			if (!response.ok) {
				// 2. Tangkap error plain text dari Backend (Zero Trust/Validation error)
				const textError = await response.text();
				throw new Error(textError || 'Gagal menyimpan data petugas');
			}

			successMessage = isEditMode
				? 'Data petugas berhasil diperbarui'
				: 'Petugas berhasil ditambahkan';

			// Beri waktu user membaca sukses sebelum refresh
			setTimeout(async () => {
				closeModal();
				await getAllUsers();
			}, 1000);
		} catch (error) {
			modalError = error.message;
		} finally {
			isLoading = false;
		}
	};

	const deleteUser = async (userId) => {
		if (!confirm('Apakah Anda yakin ingin menghapus petugas ini?')) {
			return;
		}

		try {
			const response = await apiService.delete(`/users/${userId}`);

			if (!response.ok) {
				const textError = await response.text();
				throw new Error(textError || 'Gagal menghapus petugas');
			}

			await getAllUsers();
		} catch (error) {
			tableError = error.message;
			console.error(error);
		}
	};
</script>

<div class="mx-auto max-w-7xl p-4">
	<div class="mb-8">
		<h1 class="text-3xl font-bold text-gray-800">Data Petugas</h1>
		<p class="text-gray-600">Manajemen informasi personil KUA Kecamatan Ciemas</p>
	</div>

	<div class="mb-6">
		{#if hasAllPermissions(user.permissions, ['user:create'])}
			<button class="btn btn-primary" onclick={openModal}>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="mr-1 h-5 w-5"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M12 6v6m0 0v6m0-6h6m-6 0H6"
					/>
				</svg>
				Tambah Petugas
			</button>
		{/if}
	</div>

	{#if tableError}
		<div class="mb-4 alert alert-error shadow-sm">
			<svg
				xmlns="http://www.w3.org/2000/svg"
				class="h-6 w-6 shrink-0 stroke-current"
				fill="none"
				viewBox="0 0 24 24"
				><path
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="2"
					d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
				/></svg
			>
			<span>{tableError}</span>
		</div>
	{/if}

	<div class="card border border-base-200 bg-base-100 shadow-xl">
		<div class="card-body overflow-x-auto">
			{#if isLoading && staffs.length === 0}
				<div class="flex justify-center py-10">
					<span class="loading loading-lg loading-spinner text-primary"></span>
				</div>
			{:else}
				<table class="table table-zebra">
					<thead>
						<tr class="text-base">
							<th>Nama</th>
							<th>NIP</th>
							<th>Jabatan</th>
							<th>Peran</th>
							<th>Unit Kerja</th>
							<th>Kontak</th>
							<th>Email</th>
							<th class="text-center">Aksi</th>
						</tr>
					</thead>
					<tbody>
						{#each staffs as staff (staff.id)}
							<tr>
								<td class="font-medium">{staff.nama_lengkap}</td>
								<td class="font-mono text-sm">{staff.nip}</td>
								<td><span>{staff.jabatan}</span></td>
								<td>
									<span
										class="badge border-primary/30 bg-primary/20 px-3 py-3 font-semibold text-primary backdrop-blur-md badge-primary"
									>
										{staff.peran.toUpperCase()}
									</span>
								</td>
								<td>{staff.unit_kerja}</td>
								<td>{staff.nomor_telepon}</td>
								<td>{staff.email}</td>
								<td class="text-center">
									<div class="flex justify-center gap-2">
										{#if hasAllPermissions(user.permissions, ['user:update'])}
											<button
												class="btn text-primary btn-ghost btn-xs hover:bg-primary/10"
												onclick={() => openEditModal(staff)}
											>
												Edit
											</button>
										{/if}
										{#if hasAllPermissions(user.permissions, ['user:delete'])}
											<button
												class="btn text-error btn-ghost btn-xs hover:bg-error/10"
												onclick={() => deleteUser(staff.id)}
											>
												Hapus
											</button>
										{/if}
									</div>
								</td>
							</tr>
						{:else}
							<tr>
								<td colspan="8" class="text-center py-10 text-gray-400">Belum ada data petugas.</td>
							</tr>
						{/each}
					</tbody>
				</table>
			{/if}
		</div>
	</div>

	<dialog class="modal" open={isModalOpen}>
		<div class="modal-box max-w-2xl">
			<h3 class="mb-6 border-b pb-2 text-xl font-bold">
				{#if isEditMode}Edit Data Petugas{:else}Registrasi Petugas Baru{/if}
			</h3>

			{#if modalError}
				<div class="mb-6 alert py-2 text-sm alert-error">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="h-4 w-4 shrink-0 stroke-current"
						fill="none"
						viewBox="0 0 24 24"
						><path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
						/></svg
					>
					<span>{modalError}</span>
				</div>
			{/if}

			{#if successMessage}
				<div class="mb-6 alert py-2 text-sm font-medium alert-success">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="h-4 w-4 shrink-0 stroke-current"
						fill="none"
						viewBox="0 0 24 24"
						><path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
						/></svg
					>
					<span>{successMessage}</span>
				</div>
			{/if}

			<form
				onsubmit={(e) => {
					e.preventDefault();
					handleSubmit();
				}}
			>
				<div class="grid grid-cols-1 gap-x-6 gap-y-4 md:grid-cols-2">
					<div class="form-control">
						<label class="label pt-0" for="nip"
							><span class="label-text font-semibold">NIP</span></label
						>
						<input
							id="nip"
							type="text"
							class="input-bordered input w-full"
							bind:value={formData.nip}
							placeholder="18 digit angka"
							maxlength="18"
						/>
					</div>

					<div class="form-control">
						<label class="label pt-0" for="nama"
							><span class="label-text font-semibold">Nama Lengkap</span></label
						>
						<input
							id="nama"
							type="text"
							class="input-bordered input w-full"
							bind:value={formData.nama_lengkap}
							placeholder="Nama sesuai SK"
						/>
					</div>

					<div class="form-control">
						<label class="label pt-0" for="jabatan"
							><span class="label-text font-semibold">Jabatan</span></label
						>
						<input
							id="jabatan"
							type="text"
							class="input-bordered input w-full"
							bind:value={formData.jabatan}
							placeholder="Contoh: Penghulu"
						/>
					</div>

					<div class="form-control">
						<label class="label pt-0" for="unit_kerja">
							<span class="label-text font-semibold"
								>Unit Kerja <span class="text-xs font-normal opacity-50">(Opsional)</span></span
							>
						</label>
						<input
							id="unit_kerja"
							type="text"
							class="input-bordered input w-full"
							bind:value={formData.unit_kerja}
							placeholder="Default: KUA Ciemas"
						/>
					</div>

					<div class="form-control">
						<label class="label pt-0" for="peran">
							<span class="label-text font-semibold">Peran</span>
						</label>
						<select id="peran" class="input-bordered input w-full" bind:value={formData.peran}>
							<option value="" disabled>Pilih peran</option>
							{#each roles as role (role.ID)}
								<option value={role.name}>{role.name}</option>
							{/each}
						</select>
					</div>

					<div class="form-control">
						<label class="label pt-0" for="email"
							><span class="label-text font-semibold">Email</span></label
						>
						<input
							id="email"
							type="email"
							class="input-bordered input w-full"
							bind:value={formData.email}
							placeholder="petugas@kemenag.go.id"
						/>
					</div>

					<div class="form-control">
						<label class="label pt-0" for="telp"
							><span class="label-text font-semibold">Nomor Telepon</span></label
						>
						<input
							id="telp"
							type="tel"
							class="input-bordered input w-full"
							bind:value={formData.nomor_telepon}
							placeholder="08xxxxxxxxxx"
						/>
					</div>
				</div>

				<div class="modal-action mt-8">
					<button type="button" class="btn btn-ghost" onclick={closeModal} disabled={isLoading}
						>Batal</button
					>
					<button type="submit" class="btn px-8 btn-primary" disabled={isLoading}>
						{#if isLoading}
							<span class="loading loading-sm loading-spinner"></span>
						{/if}
						{#if isEditMode}Perbarui Data{:else}Simpan Data{/if}
					</button>
				</div>
			</form>
		</div>
		<button
			aria-label="Close modal"
			class="modal-backdrop bg-black/40 backdrop-blur-sm"
			onclick={closeModal}
		></button>
	</dialog>
</div>
