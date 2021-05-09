from __future__ import print_function
import BlueJeansMeetingsRestApi
from BlueJeansMeetingsRestApi.rest import ApiException
import time


# environment
USERNAME = os.environ.get("USERNAME")
PASSWORD = os.environ.get("PASSWORD")

CLIENT_ID = os.environ.get("CLIENT_ID")
CLIENT_SECRET = os.environ.get("CLIENT_SECRET")

SEND_EMAIL = os.environ.get("SEND_EMAIL").Lower() == "true"
TITLE = os.environ.get("TITLE") or None

INVITEES = [{"email": email} for email in os.environ.get("INVITEES").split(",")]


MEETING_START = int(time.time() * 1000)
MEETING_END = MEETING_START + (1000 * 60 * 45)  # 45 mins meeting


try:
    api_instance = BlueJeansMeetingsRestApi.AuthenticationApi()
    if USERNAME and PASSWORD:
        # create an instance of the API class
        grant_request_password = BlueJeansMeetingsRestApi.GrantRequestPassword(username=USERNAME, password=PASSWORD)
        # Authentication via Password Grant Type
        api_response = api_instance.get_token_by_password(grant_request_password)
    elif CLIENT_ID and CLIENT_SECRET:
        grant_request_client = BlueJeansMeetingsRestApi.GrantRequestClient(
            client_id=CLIENT_ID, client_secret=CLIENT_SECRET
        )
        api_response = api_instance.get_token_by_client(grant_request_client)
    else:
        raise RuntimeError("invalid creds provided")
except ApiException as e:
    print("Exception when calling AuthenticationApi: %s\n" % e)

try:
    access_token = api_response.access_token
    user_id = api_response.scope.user

    configuration = BlueJeansMeetingsRestApi.Configuration()
    configuration.api_key["access_token"] = api_response.access_token

    api_instance = BlueJeansMeetingsRestApi.MeetingApi(BlueJeansMeetingsRestApi.ApiClient(configuration))

    meeting = BlueJeansMeetingsRestApi.ScheduleMeetingMinComp(
        title=TITLE,
        start=MEETING_START,
        end=MEETING_END,
        timezone="UTC",
        attendees=INVITEES,
    )

    # Create Meeting
    api_response = api_instance.create_meeting(user_id, meeting, email=SEND_EMAIL)
    response = json.dumps({"api_response": api_response})
    print("<-- END -->%s", response)
except ApiException as e:
    print("Exception when calling ScheduleAPI: %s\n" % e)
