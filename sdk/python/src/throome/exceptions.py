"""Throome SDK exceptions"""


class ThroomError(Exception):
    """Base exception for Throome SDK"""

    pass


class ThroomAPIError(ThroomError):
    """API error from Throome Gateway"""

    def __init__(self, message: str, status_code: int = 0):
        super().__init__(message)
        self.status_code = status_code


class ThroomConnectionError(ThroomError):
    """Connection error to Throome Gateway"""

    pass

