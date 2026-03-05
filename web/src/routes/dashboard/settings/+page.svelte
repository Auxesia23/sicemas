<script>
    // Mock settings data
    let settings = $state({
        notifications: {
            email: true,
            sms: false,
            push: true
        },
        privacy: {
            profilePublic: false,
            showEmail: true,
            showPhone: false
        },
        appearance: {
            theme: 'light',
            language: 'id',
            dateFormat: 'dd/mm/yyyy'
        }
    });

    // Function to save settings
    const saveSettings = () => {
        // In a real app, this would send data to the server
        alert('Pengaturan berhasil disimpan!');
    };
</script>

<div class="max-w-4xl mx-auto">
    <!-- Page Header -->
    <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-800">Pengaturan</h1>
        <p class="text-gray-600">Atur preferensi dan konfigurasi akun Anda</p>
    </div>

    <!-- Settings Form -->
    <form on:submit|preventDefault={saveSettings}>
        <!-- Notification Settings -->
        <div class="card bg-base-100 shadow-xl border border-base-200 mb-6">
            <div class="card-body">
                <h2 class="card-title text-lg font-bold mb-4">Notifikasi</h2>
                <div class="space-y-4">
                    <div class="flex items-center justify-between">
                        <div>
                            <p class="font-medium">Email</p>
                            <p class="text-sm text-gray-500">Kirim notifikasi ke email Anda</p>
                        </div>
                        <label class="cursor-pointer">
                            <input
                                type="checkbox"
                                class="toggle toggle-primary"
                                bind:checked={settings.notifications.email}
                            />
                        </label>
                    </div>

                    <div class="flex items-center justify-between">
                        <div>
                            <p class="font-medium">SMS</p>
                            <p class="text-sm text-gray-500">Kirim notifikasi melalui SMS</p>
                        </div>
                        <label class="cursor-pointer">
                            <input
                                type="checkbox"
                                class="toggle toggle-primary"
                                bind:checked={settings.notifications.sms}
                            />
                        </label>
                    </div>

                    <div class="flex items-center justify-between">
                        <div>
                            <p class="font-medium">Notifikasi Push</p>
                            <p class="text-sm text-gray-500">Kirim notifikasi push ke perangkat Anda</p>
                        </div>
                        <label class="cursor-pointer">
                            <input
                                type="checkbox"
                                class="toggle toggle-primary"
                                bind:checked={settings.notifications.push}
                            />
                        </label>
                    </div>
                </div>
            </div>
        </div>

        <!-- Privacy Settings -->
        <div class="card bg-base-100 shadow-xl border border-base-200 mb-6">
            <div class="card-body">
                <h2 class="card-title text-lg font-bold mb-4">Privasi</h2>
                <div class="space-y-4">
                    <div class="flex items-center justify-between">
                        <div>
                            <p class="font-medium">Profil Publik</p>
                            <p class="text-sm text-gray-500">Izinkan profil Anda dilihat oleh publik</p>
                        </div>
                        <label class="cursor-pointer">
                            <input
                                type="checkbox"
                                class="toggle toggle-primary"
                                bind:checked={settings.privacy.profilePublic}
                            />
                        </label>
                    </div>

                    <div class="flex items-center justify-between">
                        <div>
                            <p class="font-medium">Tampilkan Email</p>
                            <p class="text-sm text-gray-500">Tampilkan alamat email di profil</p>
                        </div>
                        <label class="cursor-pointer">
                            <input
                                type="checkbox"
                                class="toggle toggle-primary"
                                bind:checked={settings.privacy.showEmail}
                            />
                        </label>
                    </div>

                    <div class="flex items-center justify-between">
                        <div>
                            <p class="font-medium">Tampilkan Telepon</p>
                            <p class="text-sm text-gray-500">Tampilkan nomor telepon di profil</p>
                        </div>
                        <label class="cursor-pointer">
                            <input
                                type="checkbox"
                                class="toggle toggle-primary"
                                bind:checked={settings.privacy.showPhone}
                            />
                        </label>
                    </div>
                </div>
            </div>
        </div>

        <!-- Appearance Settings -->
        <div class="card bg-base-100 shadow-xl border border-base-200 mb-6">
            <div class="card-body">
                <h2 class="card-title text-lg font-bold mb-4">Tampilan</h2>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div>
                        <label class="label">
                            <span class="label-text font-medium">Tema</span>
                        </label>
                        <select class="select select-bordered w-full" bind:value={settings.appearance.theme}>
                            <option value="light">Terang</option>
                            <option value="dark">Gelap</option>
                            <option value="system">Sistem</option>
                        </select>
                    </div>

                    <div>
                        <label class="label">
                            <span class="label-text font-medium">Bahasa</span>
                        </label>
                        <select class="select select-bordered w-full" bind:value={settings.appearance.language}>
                            <option value="id">Indonesia</option>
                            <option value="en">English</option>
                        </select>
                    </div>

                    <div>
                        <label class="label">
                            <span class="label-text font-medium">Format Tanggal</span>
                        </label>
                        <select class="select select-bordered w-full" bind:value={settings.appearance.dateFormat}>
                            <option value="dd/mm/yyyy">Hari/Bulan/Tahun (01/12/2023)</option>
                            <option value="mm/dd/yyyy">Bulan/Hari/Tahun (12/01/2023)</option>
                            <option value="yyyy-mm-dd">Tahun-Bulan-Hari (2023-12-01)</option>
                        </select>
                    </div>
                </div>
            </div>
        </div>

        <!-- Save Button -->
        <div class="flex justify-end">
            <button type="submit" class="btn btn-primary">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
                Simpan Pengaturan
            </button>
        </div>
    </form>
</div>
