<script>
	import { onMount } from 'svelte';
	import { z } from 'zod';
	import apiService from '$lib/api';
	import { hasAllPermissions } from '$lib/permissions.js';
	import UserTableActions from '$lib/components/ui/UserTableActions.svelte';
	import EditUserModal from '$lib/components/ui/EditUserModal.svelte';
	import ConfirmModal from '$lib/components/ui/ConfirmModal.svelte';
	import Toast from '$lib/components/ui/Toast.svelte';

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
	let formErrors = $state({}); // State untuk menangkap error per field dari Zod

	// Modal and toast states for ConfirmModal & Toast
	let showDeleteModal = $state(false);
	let userToDelete = $state(null);
	let isDeleting = $state(false);
	let showToast = $state(false);
	let toastMessage = $state('');
	let toastType = $state('success');

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

	// Zod validation schema (Hanya butuh SATU saja)
	const userSchema = z.object({
		nip: z
			.string()
			.regex(/^(\d{16}|\d{18})$/, 'Harus berupa NIK (16 digit angka) atau NIP (18 digit angka)'),
		nama_lengkap: z.string().min(1, 'Nama lengkap wajib diisi'),
		email: z
			.string()
			.email('Format email tidak valid')
			.or(z.literal('')) // Mengizinkan string kosong jika user tidak mengisi
			.optional(),
		nomor_telepon: z
			.string()
			.min(1, 'Nomor telepon wajib diisi')
			.regex(/^[0-9+\-\s]+$/, 'Format nomor telepon tidak valid (hanya angka, +, atau -)'),
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
		formErrors = {};
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
		formErrors = {};
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
		formErrors = {};
	};

	const handleSubmit = async () => {
		modalError = '';
		formErrors = {}; // Reset error field sebelum validasi ulang

		// 1. Validasi Client-side dengan Zod (Hanya panggil SATU kali)
		const validation = userSchema.safeParse(formData);

		if (!validation.success) {
			// Mapping error Zod ke dalam object formErrors
			validation.error.issues.forEach((issue) => {
				formErrors[issue.path[0]] = issue.message;
			});
			// Stop eksekusi agar tidak memanggil API jika ada field yang salah
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
				// Tangkap error plain text dari Backend (Zero Trust/Validation error)
				const textError = await response.text();
				throw new Error(textError || 'Gagal menyimpan data petugas');
			}

			// SUCCESS: Close modal first, then show Toast
			closeModal();

			toastMessage = isEditMode
				? 'Data petugas berhasil diperbarui'
				: 'Petugas berhasil ditambahkan';
			toastType = 'success';
			showToast = true;

			await getAllUsers(); // Refresh data table

			// Hide Toast after 3 seconds
			setTimeout(() => {
				showToast = false;
			}, 3000);
		} catch (error) {
			modalError = error.message; // Keep API error inside modal
		} finally {
			isLoading = false;
		}
	};

	// Confirm delete action
	const confirmDelete = (user) => {
		userToDelete = user;
		showDeleteModal = true;
	};

	// Handle delete action with Toast
	const deleteUser = async () => {
		isDeleting = true;

		try {
			const response = await apiService.delete(`/users/${userToDelete.id}`);

			if (!response.ok) {
				const textError = await response.text();
				throw new Error(textError || 'Gagal menghapus petugas');
			}

			// Success
			showDeleteModal = false;
			toastMessage = 'Petugas berhasil dihapus';
			toastType = 'success';
			showToast = true;
			await getAllUsers();
			setTimeout(() => {
				showToast = false;
			}, 3000);
		} catch (error) {
			// Error
			tableError = error.message;
			toastMessage = error.message;
			toastType = 'error';
			showToast = true;
			setTimeout(() => {
				showToast = false;
			}, 4000);
		} finally {
			isDeleting = false;
		}
	};
</script>

<div class="mx-auto max-w-7xl p-4">
	<Toast
		show={showToast}
		message={toastMessage}
		type={toastType}
		onclose={() => (showToast = false)}
	/>

	<ConfirmModal
		show={showDeleteModal}
		title="Hapus Petugas"
		message="Apakah Anda yakin ingin menghapus {userToDelete?.nama_lengkap}? Tindakan ini tidak dapat dibatalkan."
		isProcessing={isDeleting}
		onConfirm={deleteUser}
		onCancel={() => (showDeleteModal = false)}
	/>

	<EditUserModal
		show={isModalOpen}
		{isEditMode}
		{isLoading}
		bind:formData
		{roles}
		{modalError}
		bind:errors={formErrors}
		onSubmit={handleSubmit}
		onCancel={closeModal}
	/>

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
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="2"
					d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
				/>
			</svg>
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
							<th>NIP / NIK</th>
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
								<td>
									<UserTableActions
										id={staff.id}
										permissions={user.permissions}
										onUpdate={() => openEditModal(staff)}
										onDelete={() => confirmDelete(staff)}
									/>
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
</div>
