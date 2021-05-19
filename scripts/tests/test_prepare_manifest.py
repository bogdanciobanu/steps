import unittest
import tempfile
import prepare_manifests
import ruamel.yaml

TEST_MANIFEST = """
apiVersion: stackpulse.io/v1
kind: Step
metadata:
  version: 1.0.0
  name: virustotal_get_ip
  description: Retrieve information about an IP address.
  vendor: virustotal
  type: triage
envs:
- name: IP
  description: IP address.
  type: string
  required: true
  example: 8.8.8.8

outputs:
- name: as_owner
  description: Owner of the Autonomous System to which the IP belongs.
  type: string
  example: Strato AG
- name: country
  description: Country where the IP is placed (ISO-3166 country code).
  type: string
  example: DE
- name: reputation
  description: IP score calculated from the votes of the VirusTotal's community.
  type: int
  example: '0'
- name: last_analysis_stats
  description: Number of different results from this scans.
  type: json
  example: |
    {"harmless": 5, "malicious": 0, "suspicious": 0, "timeout": 0, "undetected": 0}
- name: api_object
  description: Object containing an the response from the VirusTotal API.
  type: json
  example: |
    {
      "type": "ip_address",
      "id": "1.1.1.1",
      "links": {
        "self": "https://www.virustotal.com/api/v3/ip_addresses/1.1.1.1"
      },
      "data": {
        "attributes": {
          "as_owner": "Cloudflare Inc.",
          "asn": 13335,
          "country": "US"
        }
      }
    }

integrations:
- virustotal_token

"""

class PrepareManifestTest(unittest.TestCase):
    def test_patch_image_name_multiline(self):
        fp = tempfile.NamedTemporaryFile(delete=False)
        fp.write(TEST_MANIFEST.encode())
        fp.close()

        image_name = "us-docker.pkg.dev/stackpulse/public/virustotal/get-ip"
        prepare_manifests.patch_image_name(fp.name, image_name)

        yaml = ruamel.yaml.YAML()
        with open(fp.name) as f:
            py = yaml.load(f)

        test_manifest = yaml.load(TEST_MANIFEST)

        self.assertEqual(py['metadata']['imageName'], image_name)
        for i, o in enumerate(py['outputs']):
            self.assertNotIn("\\", o['example'])
            self.assertEqual(test_manifest['outputs'][i]['example'], o['example'])


if __name__ == '__main__':
    unittest.main()
