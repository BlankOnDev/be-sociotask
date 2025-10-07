package pages

// sementara, butuh test
const IndexPage = `
<html>
	<head>
		<title>OAuth-2 Test</title>
	</head>
	<body>
		<h2>OAuth-2 Test</h2>
		<p>
			Login with the following,
		</p>
		<ul>
			<li><a href="/login/google">Google</a></li>
		</ul>
	</body>
</html>
`

const SuccessPage = `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Login Success</title>
		<style>
			body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; background-color: #f0fdf4; color: #15803d; }
			.container { text-align: center; background-color: #fff; padding: 40px; border-radius: 12px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); border: 1px solid #bbf7d0; }
			h1 { font-size: 2em; margin-bottom: 20px; }
			p { margin-bottom: 20px; }
			.token-box { word-break: break-all; background-color: #f0fdf4; border: 1px dashed #86efac; padding: 15px; border-radius: 8px; margin-bottom: 20px; font-family: monospace; }
			button { background-color: #22c55e; color: white; border: none; padding: 10px 20px; border-radius: 8px; font-size: 1em; cursor: pointer; transition: background-color 0.2s; }
			button:hover { background-color: #16a34a; }
			.copied-msg { color: #16a34a; font-weight: bold; opacity: 0; transition: opacity 0.5s; }
		</style>
	</head>
	<body>
		<div class="container">
			<h1>✅ Login Successful!</h1>
			<p>Your JWT token is ready. Use this for authorized API calls.</p>
			<div id="token-display" class="token-box">Loading token...</div>
			<button onclick="copyToken()">Copy Token</button>
			<p id="copied" class="copied-msg">Copied!</p>
		</div>

		<script>
			// Ambil token dari URL parameter
			const params = new URLSearchParams(window.location.search);
			const token = params.get('token');
			const tokenDisplay = document.getElementById('token-display');

			if (token) {
				tokenDisplay.textContent = token;
			} else {
				tokenDisplay.textContent = 'Token not found in URL.';
			}

			function copyToken() {
				if (token) {
					navigator.clipboard.writeText(token).then(() => {
						const copiedMsg = document.getElementById('copied');
						copiedMsg.style.opacity = 1;
						setTimeout(() => { copiedMsg.style.opacity = 0; }, 2000);
					});
				}
			}
		</script>
	</body>
	</html>
`

const FailedPage = `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Login Failed</title>
		<style>
			body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; background-color: #fff1f2; color: #be123c; }
			.container { text-align: center; background-color: #fff; padding: 40px; border-radius: 12px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); border: 1px solid #fecdd3; }
			h1 { font-size: 2em; margin-bottom: 20px; }
			p { margin-bottom: 20px; }
			.error-box { background-color: #fff1f2; border: 1px dashed #fda4af; padding: 15px; border-radius: 8px; font-family: monospace; }
			a { color: #1d4ed8; text-decoration: none; font-weight: bold; }
			a:hover { text-decoration: underline; }
		</style>
	</head>
	<body>
		<div class="container">
			<h1>❌ Login Failed</h1>
			<p>Something went wrong. Here is the error detail:</p>
			<div id="error-display" class="error-box">No error reason provided.</div>
			<p style="margin-top: 30px;"><a href="/">Try Again</a></p>
		</div>

		<script>
			// Ambil pesan error dari URL parameter
			const params = new URLSearchParams(window.location.search);
			const error = params.get('error');
			const errorDisplay = document.getElementById('error-display');

			if (error) {
				errorDisplay.textContent = error.replace(/_/g, ' '); // Ganti underscore dengan spasi
			}
		</script>
	</body>
	</html>
`
