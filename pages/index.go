package pages

var Index = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Database Entries Dashboard</title>
    <link rel="stylesheet" href="/static/styles.css">
</head>
<body>
    <div class="container">
        <header>
            <h1>📊 Database Entries Dashboard</h1>
            <p>View and manage all stored records</p>
        </header>

        <section class="stats-bar">
            <div class="stat-card">
                <div class="stat-value"></div>
                <div class="stat-label">Total Entries</div>
            </div>
            <div class="stat-card">
                <div class="stat-value"></div>
                <div class="stat-label">Active</div>
            </div>
            <div class="stat-card">
                <div class="stat-value"></div>
                <div class="stat-label">Growth This Month</div>
            </div>
        </section>

        <section class="entries-table">
            <table>
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>Name</th>
                        <th>Email</th>
                        <th>Date Added</th>
                        <th>Status</th>
                        <th>Action</th>
                    </tr>
                </thead>
                <tbody>
                    <!-- Populate with dynamic data -->
                </tbody>
            </table>
        </section>

        <footer class="footer">
            <p>&copy; 2024 Your Application</p>
        </footer>
    </div>
</body>
</html>
`
