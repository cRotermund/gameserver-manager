import express from "express";

const app = express();
const port = process.env.PORT ?? 3000;

app.get("/health", (_req, res) => {
	res.json({ status: "healthy" });
});

app.listen(port, () => {
	console.log(`control-plane-web listening on :${port}`);
});
