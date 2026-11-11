# NetOps API - Documentation de l'API REST

Cette API permet d'ingérer et de consulter les métriques système transmises par les agents installés sur les machines.

- **Base URL :** `http://<server-host>:8080`
- **Format des données :** `JSON`
- **Champs temporels :** Format ISO 8601 / RFC 3339 (ex: `2026-11-09T22:29:36.2613398Z`)

---

## Format des Erreurs

Toutes les erreurs renvoyées par l'API partagent la même structure JSON :

```json
{
  "error": "Message explicatif de l'erreur",
  "code": "400"
}
```

---

## Sommaire des Endpoints

| Méthode | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/health` | Vérifie l'état de l'API et la connexion DB |
| `GET` | `/api/computers` | Liste toutes les machines enregistrées |
| `GET` | `/api/computers/{id}` | Détails d'une machine spécifique |
| `GET` | `/api/reports` | Liste l'ensemble des historiques de métriques |
| `GET` | `/api/reports/{id}` | Récupère un rapport spécifique par son ID |
| `POST` | `/api/report` | Ingestion d'un nouveau rapport (Agent) |

---

## Détails des Endpoints

### 1. Healthcheck

Vérifie que l'API est fonctionnelle et que la connexion à MariaDB est active.

* **URL :** `/health`
* **Méthode :** `GET`

#### Réponses :
* **`200 OK`** : API et base de données opérationnelles
  ```json
  {
    "status": "ok"
  }
  ```
* **`503 Service Unavailable`** : Problème de connexion avec la base de données
  ```json
  {
    "status": "ko"
  }
  ```

---

### 2. Lister les machines (`Devices`)

Récupère la liste de toutes les machines enregistrées en base.

* **URL :** `/api/devices`
* **Méthode :** `GET`

#### Réponses :
* **`200 OK`** :
  ```json
  [
    {
      "id": "6b1a3ec7-67cb-4271-aa3e-db43488b18c1",
      "hostname": "windows-dylan",
      "os": "windows",
      "platform": "Microsoft Windows 11 Home 10.0.26300.9711 Build 26300.9711",
      "arch": "amd64"
    }
  ]
  ```
* **`503 Service Unavailable`** : Erreur lors de la récupération en base de données.

---

### 3. Récupérer une machine par son ID

* **URL :** `/api/devices/{id}`
* **Méthode :** `GET`

#### Paramètres d'URL :
* `id` *(string)* : L'UUID ou l'identifiant unique de la machine.

#### Réponses :
* **`200 OK`** :
  ```json
  {
    "id": "6b1a3ec7-67cb-4271-aa3e-db43488b18c1",
    "hostname": "windows-dylan",
    "os": "windows",
    "platform": "Microsoft Windows 11 Home 10.0.26300.9711 Build 26300.9711",
    "arch": "amd64"
  }
  ```
* **`404 Not Found`** :
  ```json
  {
    "error": "machine introuvable",
    "code": "404"
  }
  ```

---

### 4. Lister les rapports d'historique (`Reports`)

Récupère l'ensemble des rapports de métriques enregistrés.

* **URL :** `/api/reports`
* **Méthode :** `GET`

#### Réponses :
* **`200 OK`** :
  ```json
  [
    {
      "id": 1,
      "send_at": "2026-11-09T22:29:36.2613398Z",
      "device_id": "6b1a3ec7-67cb-4271-aa3e-db43488b18c1",
      "ip_address": "192.168.1.186",
      "agent_version": "0.0.1-ALPHA",
      "uptime_second": 256,
      "mem_usage_percent": 30,
      "mem_used_bytes": 5243146240,
      "mem_total_bytes": 17114173440,
      "cpu_usage_percent": 8.88,
      "cpu_cores": 6
    }
  ]
  ```

---

### 5. Récupérer un rapport d'historique par son ID

* **URL :** `/api/reports/{id}`
* **Méthode :** `GET`

#### Paramètres d'URL :
* `id` *(string/int)* : L'ID de l'historique recherché.

#### Réponses :
* **`200 OK`** :
  ```json
  {
    "id": 1,
    "send_at": "2026-11-09T22:29:36.2613398Z",
    "device_id": "6b1a3ec7-67cb-4271-aa3e-db43488b18c1",
    "ip_address": "192.168.1.186",
    "agent_version": "0.0.1-ALPHA",
    "uptime_second": 256,
    "mem_usage_percent": 30,
    "mem_used_bytes": 5243146240,
    "mem_total_bytes": 17114173440,
    "cpu_usage_percent": 8.88,
    "cpu_cores": 6
  }
  ```
* **`404 Not Found`** :
  ```json
  {
    "error": "report introuvable",
    "code": "404"
  }
  ```

---

## Exemples d'intégration PHP (cURL)

### Consommer l'API avec cURL (`GET /api/computers`)
```php
$ch = curl_init('http://localhost:8080/api/computers');
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
curl_setopt($ch, CURLOPT_HTTPHEADER, ['Accept: application/json']);

$response = curl_exec($ch);
$httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
curl_close($ch);

if ($httpCode === 200) {
    $computers = json_decode($response, true);
    foreach ($computers as $computer) {
        echo "Machine : {$computer['hostname']} ({$computer['os']})\n";
    }
} else {
    echo "Erreur ($httpCode) : $response";
}
```